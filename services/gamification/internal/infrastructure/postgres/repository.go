package postgres

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"cdek/platform/gamification/internal/domain"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetState(userID string) (*domain.PortalState, error) {
	ctx := context.Background()

	state, err := r.loadState(ctx, r.pool, userID)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (r *Repository) AcceptTask(userID, taskID string) (*domain.PortalState, error) {
	ctx := context.Background()

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(
		ctx,
		`update gamification.tasks
		 set status = $3
		 where user_id = $1 and id = $2 and status = $4`,
		userID,
		taskID,
		domain.TaskStatusInProgress,
		domain.TaskStatusAvailable,
	)
	if err != nil {
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		var exists bool
		if scanErr := tx.QueryRow(
			ctx,
			`select exists(select 1 from gamification.tasks where user_id = $1 and id = $2)`,
			userID,
			taskID,
		).Scan(&exists); scanErr != nil {
			return nil, scanErr
		}
		if !exists {
			return nil, domain.ErrTaskNotFound
		}
	}

	title, err := r.getTaskTitle(ctx, tx, userID, taskID)
	if err != nil {
		return nil, err
	}

	if err := insertNotification(ctx, tx, userID, fmt.Sprintf("notification-task-accept-%s", taskID), "Задание принято", fmt.Sprintf("Задание \"%s\" добавлено в активные.", title), "success"); err != nil {
		return nil, err
	}

	state, err := r.loadState(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return state, nil
}

func (r *Repository) AdvanceTask(userID, taskID string) (*domain.PortalState, error) {
	ctx := context.Background()

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	task, err := r.getTaskForUpdate(ctx, tx, userID, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, err
	}

	if task.Status == domain.TaskStatusAvailable {
		task.Status = domain.TaskStatusInProgress
	}

	if task.Status != domain.TaskStatusCompleted {
		task.Progress++
		if task.Progress >= task.Target {
			task.Progress = task.Target
			task.Status = domain.TaskStatusCompleted
		}

		if _, err := tx.Exec(
			ctx,
			`update gamification.tasks
			 set status = $3, progress = $4
			 where user_id = $1 and id = $2`,
			userID,
			taskID,
			task.Status,
			task.Progress,
		); err != nil {
			return nil, err
		}

		if task.Status == domain.TaskStatusCompleted {
			if err := r.applyTaskCompletion(ctx, tx, userID, task); err != nil {
				return nil, err
			}
		}
	}

	state, err := r.loadState(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return state, nil
}

func (r *Repository) RedeemReward(userID, rewardID string) (*domain.PortalState, error) {
	ctx := context.Background()

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var reward domain.Reward
	err = tx.QueryRow(
		ctx,
		`select id, title, description, cost, status, category
		 from gamification.rewards
		 where user_id = $1 and id = $2
		 for update`,
		userID,
		rewardID,
	).Scan(&reward.ID, &reward.Title, &reward.Description, &reward.Cost, &reward.Status, &reward.Category)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRewardNotFound
		}
		return nil, err
	}

	var coins int32
	if err := tx.QueryRow(
		ctx,
		`select coins from gamification.user_state where user_id = $1 for update`,
		userID,
	).Scan(&coins); err != nil {
		return nil, err
	}

	if reward.Cost > coins {
		return nil, domain.ErrInsufficientBalance
	}

	if reward.Status != domain.RewardStatusRedeemed {
		if _, err := tx.Exec(
			ctx,
			`update gamification.rewards set status = $3 where user_id = $1 and id = $2`,
			userID,
			rewardID,
			domain.RewardStatusRedeemed,
		); err != nil {
			return nil, err
		}

		if _, err := tx.Exec(
			ctx,
			`update gamification.user_state set coins = coins - $2, updated_at = now() where user_id = $1`,
			userID,
			reward.Cost,
		); err != nil {
			return nil, err
		}

		purchaseID := fmt.Sprintf("purchase-%d", time.Now().UnixNano())
		if _, err := tx.Exec(
			ctx,
			`insert into gamification.purchases (id, user_id, reward_id, title, cost, redeemed_at, status)
			 values ($1, $2, $3, $4, $5, $6, $7)`,
			purchaseID,
			userID,
			reward.ID,
			reward.Title,
			reward.Cost,
			time.Now().Format("2006-01-02"),
			"completed",
		); err != nil {
			return nil, err
		}

		if err := insertNotification(ctx, tx, userID, fmt.Sprintf("notification-reward-%s", rewardID), "Награда оформлена", fmt.Sprintf("Ты обменял баллы на \"%s\".", reward.Title), "info"); err != nil {
			return nil, err
		}
	}

	state, err := r.loadState(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return state, nil
}

func (r *Repository) applyTaskCompletion(ctx context.Context, tx pgx.Tx, userID string, task *domain.Task) error {
	if _, err := tx.Exec(
		ctx,
		`update gamification.user_state
		 set current_xp = current_xp + $2,
		     coins = coins + ($2 * 4),
		     today_earned = today_earned + $2,
		     completed_tasks = completed_tasks + 1,
		     updated_at = now()
		 where user_id = $1`,
		userID,
		task.RewardXP,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		ctx,
		`insert into gamification.recent_activity (id, user_id, title, timestamp_label, xp, sort_order)
		 values ($1, $2, $3, $4, $5, $6)`,
		fmt.Sprintf("activity-task-%d", time.Now().UnixNano()),
		userID,
		fmt.Sprintf("Закрыто задание \"%s\"", task.Title),
		"Только что",
		task.RewardXP,
		time.Now().Unix(),
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		ctx,
		`update gamification.leaderboard
		 set xp = (select current_xp from gamification.user_state where user_id = $1)
		 where user_id = $1`,
		userID,
	); err != nil {
		return err
	}

	if err := rebuildLeaderboardRanks(ctx, tx); err != nil {
		return err
	}

	if err := insertNotification(
		ctx,
		tx,
		userID,
		fmt.Sprintf("notification-task-complete-%s", task.ID),
		"Задание выполнено",
		fmt.Sprintf("Ты закрыл \"%s\" и получил %d XP.", task.Title, task.RewardXP),
		"success",
	); err != nil {
		return err
	}

	return r.unlockGuardianIfNeeded(ctx, tx, userID)
}

func (r *Repository) unlockGuardianIfNeeded(ctx context.Context, tx pgx.Tx, userID string) error {
	var completedTasks int32
	if err := tx.QueryRow(ctx, `select completed_tasks from gamification.user_state where user_id = $1`, userID).Scan(&completedTasks); err != nil {
		return err
	}

	if completedTasks < 3 {
		return nil
	}

	var status string
	var rewardXP int32
	err := tx.QueryRow(
		ctx,
		`select status, reward_xp
		 from gamification.achievements
		 where user_id = $1 and title = 'Страж'
		 for update`,
		userID,
	).Scan(&status, &rewardXP)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	if status == domain.AchievementStatusUnlocked {
		return nil
	}

	if _, err := tx.Exec(
		ctx,
		`update gamification.achievements
		 set status = $2
		 where user_id = $1 and title = 'Страж'`,
		userID,
		domain.AchievementStatusUnlocked,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		ctx,
		`update gamification.user_state
		 set current_xp = current_xp + $2,
		     coins = coins + ($2 * 3),
		     today_earned = today_earned + $2,
		     updated_at = now()
		 where user_id = $1`,
		userID,
		rewardXP,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		ctx,
		`update gamification.leaderboard
		 set xp = (select current_xp from gamification.user_state where user_id = $1)
		 where user_id = $1`,
		userID,
	); err != nil {
		return err
	}

	if err := rebuildLeaderboardRanks(ctx, tx); err != nil {
		return err
	}

	return insertNotification(ctx, tx, userID, fmt.Sprintf("notification-achievement-guardian-%s", userID), "Получено достижение", "Открыт эпический бейдж \"Страж\".", "success")
}

func (r *Repository) getTaskForUpdate(ctx context.Context, tx pgx.Tx, userID, taskID string) (*domain.Task, error) {
	task := &domain.Task{}
	err := tx.QueryRow(
		ctx,
		`select id, title, description, status, progress, target, reward_xp
		 from gamification.tasks
		 where user_id = $1 and id = $2
		 for update`,
		userID,
		taskID,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Progress, &task.Target, &task.RewardXP)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *Repository) getTaskTitle(ctx context.Context, tx pgx.Tx, userID, taskID string) (string, error) {
	var title string
	err := tx.QueryRow(
		ctx,
		`select title from gamification.tasks where user_id = $1 and id = $2`,
		userID,
		taskID,
	).Scan(&title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrTaskNotFound
		}
		return "", err
	}

	return title, nil
}

type queryer interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

func (r *Repository) loadState(ctx context.Context, db queryer, userID string) (*domain.PortalState, error) {
	state := &domain.PortalState{UserID: userID}

	if err := db.QueryRow(
		ctx,
		`select us.current_xp, us.coins, us.today_earned, us.streak_days,
		        coalesce(lb.rank, 1), us.completed_tasks, us.api_requests, us.articles_count, us.comments_count
		 from gamification.user_state us
		 left join gamification.leaderboard lb on lb.user_id = us.user_id
		 where us.user_id = $1`,
		userID,
	).Scan(
		&state.CurrentXP,
		&state.Coins,
		&state.TodayEarned,
		&state.StreakDays,
		&state.Rank,
		&state.CompletedTasks,
		&state.ApiRequests,
		&state.ArticlesCount,
		&state.CommentsCount,
	); err != nil {
		return nil, err
	}

	var err error
	if state.WeeklyActivity, err = loadWeeklyActivity(ctx, db, userID); err != nil {
		return nil, err
	}
	if state.Articles, err = loadArticleCards(ctx, db, userID); err != nil {
		return nil, err
	}
	if state.RecentActivity, err = loadRecentActivity(ctx, db, userID); err != nil {
		return nil, err
	}
	if state.Tasks, err = loadTasks(ctx, db, userID); err != nil {
		return nil, err
	}
	if state.Achievements, err = loadAchievements(ctx, db, userID); err != nil {
		return nil, err
	}
	if state.Leaderboard, err = loadLeaderboard(ctx, db); err != nil {
		return nil, err
	}
	if state.Rewards, err = loadRewards(ctx, db, userID); err != nil {
		return nil, err
	}
	if state.Purchases, err = loadPurchases(ctx, db, userID); err != nil {
		return nil, err
	}
	if state.Notifications, err = loadNotifications(ctx, db, userID); err != nil {
		return nil, err
	}

	return state, nil
}

func loadWeeklyActivity(ctx context.Context, db queryer, userID string) ([]domain.ActivityPoint, error) {
	rows, err := db.Query(ctx, `select day_code, xp from gamification.weekly_activity where user_id = $1 order by sort_order`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ActivityPoint
	for rows.Next() {
		var item domain.ActivityPoint
		if err := rows.Scan(&item.Day, &item.XP); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func loadArticleCards(ctx context.Context, db queryer, userID string) ([]domain.ArticleCard, error) {
	rows, err := db.Query(ctx, `select id, title, views, comments, xp, rating from gamification.article_cards where user_id = $1 order by title`, userID)
	if err != nil {
		if strings.Contains(err.Error(), `relation "gamification.article_cards" does not exist`) {
			return []domain.ArticleCard{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var result []domain.ArticleCard
	for rows.Next() {
		var item domain.ArticleCard
		if err := rows.Scan(&item.ID, &item.Title, &item.Views, &item.Comments, &item.XP, &item.Rating); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func loadRecentActivity(ctx context.Context, db queryer, userID string) ([]domain.ActivityItem, error) {
	rows, err := db.Query(
		ctx,
		`select id, title, timestamp_label, xp
		 from gamification.recent_activity
		 where user_id = $1
		 order by sort_order desc, id desc`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ActivityItem
	for rows.Next() {
		var item domain.ActivityItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Timestamp, &item.XP); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func loadTasks(ctx context.Context, db queryer, userID string) ([]domain.Task, error) {
	rows, err := db.Query(
		ctx,
		`select id, title, description, status, progress, target, reward_xp
		 from gamification.tasks
		 where user_id = $1
		 order by title`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Task
	for rows.Next() {
		var item domain.Task
		if err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Status, &item.Progress, &item.Target, &item.RewardXP); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func loadAchievements(ctx context.Context, db queryer, userID string) ([]domain.Achievement, error) {
	rows, err := db.Query(
		ctx,
		`select id, title, description, rarity, status, reward_xp
		 from gamification.achievements
		 where user_id = $1
		 order by title`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Achievement
	for rows.Next() {
		var item domain.Achievement
		if err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Rarity, &item.Status, &item.RewardXP); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func loadLeaderboard(ctx context.Context, db queryer) ([]domain.LeaderboardEntry, error) {
	rows, err := db.Query(ctx, `select user_id, rank, xp from gamification.leaderboard order by rank, xp desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.LeaderboardEntry
	for rows.Next() {
		var item domain.LeaderboardEntry
		if err := rows.Scan(&item.UserID, &item.Rank, &item.XP); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	sort.Slice(result, func(left, right int) bool {
		return result[left].Rank < result[right].Rank
	})

	return result, rows.Err()
}

func loadRewards(ctx context.Context, db queryer, userID string) ([]domain.Reward, error) {
	rows, err := db.Query(
		ctx,
		`select id, title, description, cost, status, category
		 from gamification.rewards
		 where user_id = $1
		 order by cost`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Reward
	for rows.Next() {
		var item domain.Reward
		if err := rows.Scan(&item.ID, &item.Title, &item.Description, &item.Cost, &item.Status, &item.Category); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func loadPurchases(ctx context.Context, db queryer, userID string) ([]domain.Purchase, error) {
	rows, err := db.Query(
		ctx,
		`select id, reward_id, title, cost, redeemed_at, status
		 from gamification.purchases
		 where user_id = $1
		 order by redeemed_at desc, id desc`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Purchase
	for rows.Next() {
		var item domain.Purchase
		if err := rows.Scan(&item.ID, &item.RewardID, &item.Title, &item.Cost, &item.RedeemedAt, &item.Status); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func loadNotifications(ctx context.Context, db queryer, userID string) ([]domain.Notification, error) {
	rows, err := db.Query(
		ctx,
		`select id, title, body, variant
		 from gamification.notifications
		 where user_id = $1
		 order by created_at desc, id desc`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Notification
	for rows.Next() {
		var item domain.Notification
		if err := rows.Scan(&item.ID, &item.Title, &item.Body, &item.Variant); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, rows.Err()
}

func insertNotification(ctx context.Context, tx pgx.Tx, userID, id, title, body, variant string) error {
	_, err := tx.Exec(
		ctx,
		`insert into gamification.notifications (id, user_id, title, body, variant)
		 values ($1, $2, $3, $4, $5)
		 on conflict (id) do nothing`,
		id,
		userID,
		title,
		body,
		variant,
	)
	return err
}

func rebuildLeaderboardRanks(ctx context.Context, tx pgx.Tx) error {
	rows, err := tx.Query(ctx, `select user_id from gamification.leaderboard order by xp desc, user_id asc`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return err
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for index, userID := range userIDs {
		if _, err := tx.Exec(
			ctx,
			`update gamification.leaderboard set rank = $2 where user_id = $1`,
			userID,
			index+1,
		); err != nil {
			return err
		}
	}

	return nil
}
