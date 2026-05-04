package app

import (
	"context"
	"log/slog"
	"sort"

	gamificationv1 "cdek/platform/shared/contracts/gamification/v1"
	userv1 "cdek/platform/shared/contracts/user/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	logger      *slog.Logger
	tracer      trace.Tracer
	db          *pgxpool.Pool
	userClient  userv1.UserServiceClient
	gameClient  gamificationv1.GamificationServiceClient
}

func NewServer(
	logger *slog.Logger,
	tracer trace.Tracer,
	db *pgxpool.Pool,
	userClient userv1.UserServiceClient,
	gameClient gamificationv1.GamificationServiceClient,
) *Server {
	return &Server{
		logger:     logger,
		tracer:     tracer,
		db:         db,
		userClient: userClient,
		gameClient: gameClient,
	}
}

func (s *Server) Health(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": "ok",
	})
}

func (s *Server) Bootstrap(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	requestContext, span := s.tracer.Start(ctx.UserContext(), "bootstrap")
	defer span.End()

	payload, err := s.composePayload(requestContext, userID, nil)
	if err != nil {
		s.logger.Error("bootstrap failed", "error", err)
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.JSON(payload)
}

func (s *Server) AcceptTask(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	requestContext, span := s.tracer.Start(ctx.UserContext(), "accept-task")
	defer span.End()

	response, err := s.gameClient.AcceptTask(requestContext, &gamificationv1.AcceptTaskRequest{
		UserId: userID,
		TaskId: ctx.Params("taskId"),
	})
	if err != nil {
		s.logger.Error("accept task failed", "error", err)
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": err.Error()})
	}

	payload, err := s.composePayload(requestContext, userID, response.Snapshot)
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.JSON(payload)
}

func (s *Server) AdvanceTask(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	requestContext, span := s.tracer.Start(ctx.UserContext(), "advance-task")
	defer span.End()

	response, err := s.gameClient.AdvanceTask(requestContext, &gamificationv1.AdvanceTaskRequest{
		UserId: userID,
		TaskId: ctx.Params("taskId"),
	})
	if err != nil {
		s.logger.Error("advance task failed", "error", err)
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": err.Error()})
	}

	payload, err := s.composePayload(requestContext, userID, response.Snapshot)
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.JSON(payload)
}

func (s *Server) RedeemReward(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	requestContext, span := s.tracer.Start(ctx.UserContext(), "redeem-reward")
	defer span.End()

	response, err := s.gameClient.RedeemReward(requestContext, &gamificationv1.RedeemRewardRequest{
		UserId:   userID,
		RewardId: ctx.Params("rewardId"),
	})
	if err != nil {
		s.logger.Error("redeem reward failed", "error", err)
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": err.Error()})
	}

	payload, err := s.composePayload(requestContext, userID, response.Snapshot)
	if err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.JSON(payload)
}

func (s *Server) composePayload(ctx context.Context, userID string, snapshot *gamificationv1.GetPortalSnapshotResponse) (fiber.Map, error) {
	currentUser, err := s.userClient.GetCurrentUser(ctx, &userv1.GetCurrentUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}

	if snapshot == nil {
		snapshot, err = s.gameClient.GetPortalSnapshot(ctx, &gamificationv1.GetPortalSnapshotRequest{UserId: userID})
		if err != nil {
			return nil, err
		}
	}

	directory, err := s.loadDirectory(ctx, currentUser.User, snapshot.Leaderboard)
	if err != nil {
		return nil, err
	}

	return fiber.Map{
		"currentUser":   composeCurrentUser(currentUser.User, snapshot.Summary),
		"dashboard":     composeDashboard(snapshot),
		"tasks":         composeTasks(snapshot.Tasks),
		"achievements":  composeAchievements(snapshot),
		"leaderboard":   composeLeaderboard(snapshot.Leaderboard, directory, userID),
		"rewards":       composeRewards(snapshot.Rewards),
		"purchases":     composePurchases(snapshot.Purchases),
		"notifications": composeNotifications(snapshot.Notifications),
	}, nil
}

func (s *Server) loadDirectory(ctx context.Context, currentUser *userv1.UserSummary, leaderboard []*gamificationv1.LeaderboardEntry) (map[string]*userv1.UserSummary, error) {
	userIDs := []string{currentUser.Id}
	for _, entry := range leaderboard {
		userIDs = append(userIDs, entry.UserId)
	}

	response, err := s.userClient.BatchGetUsers(ctx, &userv1.BatchGetUsersRequest{UserIds: unique(userIDs)})
	if err != nil {
		return nil, err
	}

	directory := make(map[string]*userv1.UserSummary, len(response.Users)+1)
	directory[currentUser.Id] = currentUser
	for _, user := range response.Users {
		directory[user.Id] = user
	}

	return directory, nil
}

func unique(items []string) []string {
	index := map[string]struct{}{}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if _, exists := index[item]; exists {
			continue
		}
		index[item] = struct{}{}
		result = append(result, item)
	}

	sort.Strings(result)
	return result
}

func composeCurrentUser(user *userv1.UserSummary, summary *gamificationv1.UserSummary) fiber.Map {
	return fiber.Map{
		"id":              user.Id,
		"name":            user.Name,
		"title":           user.Title,
		"company":         user.Company,
		"joinedAt":        user.JoinedAt,
		"location":        user.Location,
		"team":            user.Team,
		"level":           summary.Level,
		"levelText":       summary.LevelText,
		"currentXp":       summary.CurrentXp,
		"progressPercent": summary.ProgressPercent,
		"xpToNextLevel":   summary.XpToNextLevel,
		"coins":           summary.Coins,
		"todayEarned":     summary.TodayEarned,
		"streakDays":      summary.StreakDays,
		"rank":            summary.Rank,
		"totalBadges":     summary.TotalBadges,
		"completedTasks":  summary.CompletedTasks,
	}
}

func composeDashboard(snapshot *gamificationv1.GetPortalSnapshotResponse) fiber.Map {
	return fiber.Map{
		"metrics":        snapshot.Metrics,
		"weeklyActivity": snapshot.WeeklyActivity,
		"articles":       snapshot.Articles,
		"recentActivity": snapshot.RecentActivity,
	}
}

func composeTasks(tasks []*gamificationv1.TaskItem) []fiber.Map {
	result := make([]fiber.Map, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, fiber.Map{
			"id":          task.Id,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"progress":    task.Progress,
			"target":      task.Target,
			"rewardXp":    task.RewardXp,
		})
	}

	return result
}

func composeAchievements(snapshot *gamificationv1.GetPortalSnapshotResponse) fiber.Map {
	return fiber.Map{
		"buckets": snapshot.AchievementBuckets,
		"items":   snapshot.Achievements,
	}
}

func composeLeaderboard(entries []*gamificationv1.LeaderboardEntry, directory map[string]*userv1.UserSummary, currentUserID string) fiber.Map {
	rows := make([]fiber.Map, 0, len(entries))
	for _, entry := range entries {
		user := directory[entry.UserId]
		if user == nil {
			continue
		}

		rows = append(rows, fiber.Map{
			"userId":      entry.UserId,
			"rank":        entry.Rank,
			"xp":          entry.Xp,
			"levelText":   entry.LevelText,
			"name":        user.Name,
			"title":       user.Title,
			"company":     user.Company,
			"isCurrent":   entry.UserId == currentUserID,
		})
	}

	sort.Slice(rows, func(left, right int) bool {
		return rows[left]["rank"].(int32) < rows[right]["rank"].(int32)
	})

	podium := make([]fiber.Map, 0, 3)
	for _, row := range rows {
		if len(podium) == 3 {
			break
		}
		if row["rank"].(int32) <= 3 {
			podium = append(podium, row)
		}
	}

	return fiber.Map{
		"podium": podium,
		"rows":   rows,
	}
}

func composeRewards(rewards []*gamificationv1.RewardItem) []fiber.Map {
	result := make([]fiber.Map, 0, len(rewards))
	for _, reward := range rewards {
		result = append(result, fiber.Map{
			"id":          reward.Id,
			"title":       reward.Title,
			"description": reward.Description,
			"cost":        reward.Cost,
			"status":      reward.Status,
			"category":    reward.Category,
		})
	}

	return result
}

func composePurchases(purchases []*gamificationv1.PurchaseItem) []fiber.Map {
	result := make([]fiber.Map, 0, len(purchases))
	for _, purchase := range purchases {
		result = append(result, fiber.Map{
			"id":         purchase.Id,
			"rewardId":   purchase.RewardId,
			"title":      purchase.Title,
			"cost":       purchase.Cost,
			"redeemedAt": purchase.RedeemedAt,
			"status":     purchase.Status,
		})
	}

	return result
}

func composeNotifications(notifications []*gamificationv1.NotificationItem) []fiber.Map {
	result := make([]fiber.Map, 0, len(notifications))
	for _, notification := range notifications {
		result = append(result, fiber.Map{
			"id":      notification.Id,
			"title":   notification.Title,
			"body":    notification.Body,
			"variant": notification.Variant,
		})
	}

	return result
}
