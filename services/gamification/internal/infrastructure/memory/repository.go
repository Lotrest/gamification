package memory

import (
	"fmt"
	"sync"
	"time"

	"cdek/platform/gamification/internal/domain"
)

type Repository struct {
	mu    sync.RWMutex
	state *domain.PortalState
}

func NewRepository() *Repository {
	return &Repository{
		state: &domain.PortalState{
			UserID:         "me",
			CurrentXP:      4260,
			Coins:          3850,
			TodayEarned:    380,
			StreakDays:     7,
			Rank:           12,
			CompletedTasks: 1,
			ApiRequests:    148,
			ArticlesCount:  12,
			CommentsCount:  39,
			WeeklyActivity: []domain.ActivityPoint{
				{Day: "Пн", XP: 120},
				{Day: "Вт", XP: 160},
				{Day: "Ср", XP: 110},
				{Day: "Чт", XP: 210},
				{Day: "Пт", XP: 180},
				{Day: "Сб", XP: 90},
				{Day: "Вс", XP: 130},
			},
			Articles: []domain.ArticleCard{
				{ID: "article-1", Title: "Интеграция с API CDEK: Полный гайд", Views: 2340, Comments: 18, XP: 200, Rating: "4.9"},
				{ID: "article-2", Title: "Как ускорить first paint в React SPA", Views: 1890, Comments: 11, XP: 120, Rating: "4.8"},
			},
			RecentActivity: []domain.ActivityItem{
				{ID: "activity-1", Title: "Получен бэйдж \"Детектив логов\"", Timestamp: "2 часа назад", XP: 90},
				{ID: "activity-2", Title: "Закрыто задание \"Корректор\"", Timestamp: "Сегодня", XP: 15},
				{ID: "activity-3", Title: "Опубликована статья в базе знаний", Timestamp: "Вчера", XP: 120},
			},
			Tasks: []domain.Task{
				{ID: "task-corrector", Title: "Корректор", Description: "Найти и исправить опечатку или неточность в тексте на Портале.", Status: domain.TaskStatusCompleted, Progress: 1, Target: 1, RewardXP: 15},
				{ID: "task-user-journey", Title: "В шкуре юзера", Description: "Пройти путь регистрации на Портале и оставить 1 запись.", Status: domain.TaskStatusInProgress, Progress: 1, Target: 2, RewardXP: 30},
				{ID: "task-scouting", Title: "Скаутинг", Description: "Прочитать 3 последних комментария коллег и отметить их как полезные.", Status: domain.TaskStatusInProgress, Progress: 2, Target: 3, RewardXP: 40},
				{ID: "task-reviewer", Title: "Ревизор", Description: "Пометить в коде один устаревший комментарий или метод как deprecated.", Status: domain.TaskStatusAvailable, Progress: 0, Target: 1, RewardXP: 20},
			},
			Achievements: []domain.Achievement{
				{ID: "achievement-eye", Title: "Зоркий глаз", Description: "Нашёл и исправил ошибку в контенте", Rarity: "Обычные", Status: domain.AchievementStatusUnlocked, RewardXP: 20},
				{ID: "achievement-empathy", Title: "Эмпат", Description: "Дал полезный ответ коллеге", Rarity: "Обычные", Status: domain.AchievementStatusUnlocked, RewardXP: 10},
				{ID: "achievement-listener", Title: "Слухач", Description: "Собрал качественную обратную связь", Rarity: "Обычные", Status: domain.AchievementStatusUnlocked, RewardXP: 30},
				{ID: "achievement-guardian", Title: "Страж", Description: "Закрыть 3 задания подряд без просрочки", Rarity: "Эпические", Status: domain.AchievementStatusLocked, RewardXP: 150},
			},
			Leaderboard: []domain.LeaderboardEntry{
				{UserID: "user-1", Rank: 1, XP: 87500},
				{UserID: "user-2", Rank: 2, XP: 72300},
				{UserID: "user-3", Rank: 3, XP: 48750},
				{UserID: "user-4", Rank: 4, XP: 58900},
				{UserID: "me", Rank: 12, XP: 4260},
			},
			Rewards: []domain.Reward{
				{ID: "reward-hoodie", Title: "Фирменный худи", Description: "Мерч команды платформы", Cost: 2500, Status: domain.RewardStatusAvailable, Category: "мерч"},
				{ID: "reward-coffee", Title: "Сертификат на кофе", Description: "Купон на 5 напитков", Cost: 850, Status: domain.RewardStatusAvailable, Category: "бонус"},
				{ID: "reward-dayoff", Title: "Доп. выходной", Description: "Один дополнительный day off", Cost: 3000, Status: domain.RewardStatusAvailable, Category: "time off"},
			},
			Purchases: []domain.Purchase{
				{ID: "purchase-1", RewardID: "reward-lunch", Title: "Ланч с командой", Cost: 450, RedeemedAt: "2026-04-02", Status: "completed"},
			},
			Notifications: []domain.Notification{
				{ID: "notification-1", Title: "Новый сезон стартовал", Body: "Забери первые задания недели и поднимись в рейтинге.", Variant: "info"},
				{ID: "notification-2", Title: "До нового уровня осталось 1 140 XP", Body: "Закрой одно активное задание и получи быстрый прогресс.", Variant: "success"},
			},
		},
	}
}

func (r *Repository) GetState(_ string) (*domain.PortalState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return cloneState(r.state), nil
}

func (r *Repository) AcceptTask(_ string, taskID string) (*domain.PortalState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index := range r.state.Tasks {
		if r.state.Tasks[index].ID != taskID {
			continue
		}

		if r.state.Tasks[index].Status == domain.TaskStatusAvailable {
			r.state.Tasks[index].Status = domain.TaskStatusInProgress
			r.appendNotification("notification-task-accept", "Задание принято", fmt.Sprintf("Задание \"%s\" добавлено в активные.", r.state.Tasks[index].Title), "success")
		}

		return cloneState(r.state), nil
	}

	return nil, domain.ErrTaskNotFound
}

func (r *Repository) AdvanceTask(_ string, taskID string) (*domain.PortalState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index := range r.state.Tasks {
		task := &r.state.Tasks[index]
		if task.ID != taskID {
			continue
		}

		if task.Status == domain.TaskStatusAvailable {
			task.Status = domain.TaskStatusInProgress
		}

		if task.Status == domain.TaskStatusCompleted {
			return cloneState(r.state), nil
		}

		task.Progress++
		if task.Progress >= task.Target {
			task.Progress = task.Target
			task.Status = domain.TaskStatusCompleted
			r.state.CurrentXP += task.RewardXP
			r.state.Coins += task.RewardXP * 4
			r.state.TodayEarned += task.RewardXP
			r.state.CompletedTasks++
			r.appendActivity(task)
			r.syncLeaderboard()
			r.unlockAchievementsIfNeeded()
			r.appendNotification("notification-task-complete-"+task.ID, "Задание выполнено", fmt.Sprintf("Ты закрыл \"%s\" и получил %d XP.", task.Title, task.RewardXP), "success")
		}

		return cloneState(r.state), nil
	}

	return nil, domain.ErrTaskNotFound
}

func (r *Repository) RedeemReward(_ string, rewardID string) (*domain.PortalState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index := range r.state.Rewards {
		reward := &r.state.Rewards[index]
		if reward.ID != rewardID {
			continue
		}

		if reward.Cost > r.state.Coins {
			return nil, domain.ErrInsufficientBalance
		}

		r.state.Coins -= reward.Cost
		reward.Status = domain.RewardStatusRedeemed
		r.state.Purchases = append([]domain.Purchase{
			{
				ID:         fmt.Sprintf("purchase-%d", len(r.state.Purchases)+1),
				RewardID:   reward.ID,
				Title:      reward.Title,
				Cost:       reward.Cost,
				RedeemedAt: time.Now().Format("2006-01-02"),
				Status:     "completed",
			},
		}, r.state.Purchases...)
		r.appendNotification("notification-reward-"+rewardID, "Награда оформлена", fmt.Sprintf("Ты обменял баллы на \"%s\".", reward.Title), "info")

		return cloneState(r.state), nil
	}

	return nil, domain.ErrRewardNotFound
}

func (r *Repository) appendActivity(task *domain.Task) {
	r.state.RecentActivity = append([]domain.ActivityItem{
		{
			ID:        fmt.Sprintf("activity-task-%s", task.ID),
			Title:     fmt.Sprintf("Закрыто задание \"%s\"", task.Title),
			Timestamp: "Только что",
			XP:        task.RewardXP,
		},
	}, r.state.RecentActivity...)
}

func (r *Repository) appendNotification(id, title, body, variant string) {
	r.state.Notifications = append([]domain.Notification{
		{ID: id, Title: title, Body: body, Variant: variant},
	}, r.state.Notifications...)
}

func (r *Repository) syncLeaderboard() {
	for index := range r.state.Leaderboard {
		if r.state.Leaderboard[index].UserID == "me" {
			r.state.Leaderboard[index].XP = r.state.CurrentXP
		}
	}
}

func (r *Repository) unlockAchievementsIfNeeded() {
	for index := range r.state.Achievements {
		achievement := &r.state.Achievements[index]
		if achievement.ID == "achievement-guardian" && r.state.CompletedTasks >= 3 && achievement.Status == domain.AchievementStatusLocked {
			achievement.Status = domain.AchievementStatusUnlocked
			r.state.CurrentXP += achievement.RewardXP
			r.state.TodayEarned += achievement.RewardXP
			r.state.Coins += achievement.RewardXP * 3
			r.appendNotification("notification-achievement-guardian", "Получено достижение", "Открыт эпический бэйдж \"Страж\".", "success")
		}
	}
}

func cloneState(state *domain.PortalState) *domain.PortalState {
	clone := *state
	clone.WeeklyActivity = append([]domain.ActivityPoint(nil), state.WeeklyActivity...)
	clone.Articles = append([]domain.ArticleCard(nil), state.Articles...)
	clone.RecentActivity = append([]domain.ActivityItem(nil), state.RecentActivity...)
	clone.Tasks = append([]domain.Task(nil), state.Tasks...)
	clone.Achievements = append([]domain.Achievement(nil), state.Achievements...)
	clone.Leaderboard = append([]domain.LeaderboardEntry(nil), state.Leaderboard...)
	clone.Rewards = append([]domain.Reward(nil), state.Rewards...)
	clone.Purchases = append([]domain.Purchase(nil), state.Purchases...)
	clone.Notifications = append([]domain.Notification(nil), state.Notifications...)

	return &clone
}
