package grpcserver

import (
	"context"

	gamificationv1 "cdek/platform/shared/contracts/gamification/v1"
	"cdek/platform/gamification/internal/application"
	"cdek/platform/gamification/internal/domain"
)

type Server struct {
	service *application.Service
}

func New(service *application.Service) *Server {
	return &Server{service: service}
}

func (s *Server) GetPortalSnapshot(_ context.Context, request *gamificationv1.GetPortalSnapshotRequest) (*gamificationv1.GetPortalSnapshotResponse, error) {
	snapshot, err := s.service.GetSnapshot(request.UserId)
	if err != nil {
		return nil, err
	}

	return mapSnapshot(snapshot), nil
}

func (s *Server) AcceptTask(_ context.Context, request *gamificationv1.AcceptTaskRequest) (*gamificationv1.MutationResponse, error) {
	snapshot, err := s.service.AcceptTask(request.UserId, request.TaskId)
	if err != nil {
		return nil, err
	}

	return &gamificationv1.MutationResponse{Snapshot: mapSnapshot(snapshot)}, nil
}

func (s *Server) AdvanceTask(_ context.Context, request *gamificationv1.AdvanceTaskRequest) (*gamificationv1.MutationResponse, error) {
	snapshot, err := s.service.AdvanceTask(request.UserId, request.TaskId)
	if err != nil {
		return nil, err
	}

	return &gamificationv1.MutationResponse{Snapshot: mapSnapshot(snapshot)}, nil
}

func (s *Server) RedeemReward(_ context.Context, request *gamificationv1.RedeemRewardRequest) (*gamificationv1.MutationResponse, error) {
	snapshot, err := s.service.RedeemReward(request.UserId, request.RewardId)
	if err != nil {
		return nil, err
	}

	return &gamificationv1.MutationResponse{Snapshot: mapSnapshot(snapshot)}, nil
}

func mapSnapshot(snapshot *domain.Snapshot) *gamificationv1.GetPortalSnapshotResponse {
	response := &gamificationv1.GetPortalSnapshotResponse{
		Summary: &gamificationv1.UserSummary{
			UserId:          snapshot.Summary.UserID,
			CurrentXp:       snapshot.Summary.CurrentXP,
			Level:           snapshot.Summary.Level,
			LevelText:       snapshot.Summary.LevelText,
			ProgressPercent: snapshot.Summary.ProgressPercent,
			XpToNextLevel:   snapshot.Summary.XpToNextLevel,
			Coins:           snapshot.Summary.Coins,
			TodayEarned:     snapshot.Summary.TodayEarned,
			StreakDays:      snapshot.Summary.StreakDays,
			Rank:            snapshot.Summary.Rank,
			TotalBadges:     snapshot.Summary.TotalBadges,
			CompletedTasks:  snapshot.Summary.CompletedTasks,
		},
	}

	for _, metric := range snapshot.Metrics {
		response.Metrics = append(response.Metrics, &gamificationv1.MetricCard{
			Id:      metric.ID,
			Label:   metric.Label,
			Value:   metric.Value,
			Caption: metric.Caption,
		})
	}

	for _, point := range snapshot.WeeklyActivity {
		response.WeeklyActivity = append(response.WeeklyActivity, &gamificationv1.ActivityPoint{
			Day: point.Day,
			Xp:  point.XP,
		})
	}

	for _, article := range snapshot.Articles {
		response.Articles = append(response.Articles, &gamificationv1.ArticleCard{
			Id:       article.ID,
			Title:    article.Title,
			Views:    article.Views,
			Comments: article.Comments,
			Xp:       article.XP,
			Rating:   article.Rating,
		})
	}

	for _, activity := range snapshot.RecentActivity {
		response.RecentActivity = append(response.RecentActivity, &gamificationv1.ActivityItem{
			Id:        activity.ID,
			Title:     activity.Title,
			Timestamp: activity.Timestamp,
			Xp:        activity.XP,
		})
	}

	for _, task := range snapshot.Tasks {
		response.Tasks = append(response.Tasks, &gamificationv1.TaskItem{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Progress:    task.Progress,
			Target:      task.Target,
			RewardXp:    task.RewardXP,
		})
	}

	for _, bucket := range snapshot.AchievementBuckets {
		response.AchievementBuckets = append(response.AchievementBuckets, &gamificationv1.AchievementBucket{
			Id:        bucket.ID,
			Label:     bucket.Label,
			Collected: bucket.Collected,
			Total:     bucket.Total,
			Accent:    bucket.Accent,
		})
	}

	for _, achievement := range snapshot.Achievements {
		response.Achievements = append(response.Achievements, &gamificationv1.AchievementItem{
			Id:          achievement.ID,
			Title:       achievement.Title,
			Description: achievement.Description,
			Rarity:      achievement.Rarity,
			Status:      achievement.Status,
			RewardXp:    achievement.RewardXP,
		})
	}

	for _, entry := range snapshot.Leaderboard {
		response.Leaderboard = append(response.Leaderboard, &gamificationv1.LeaderboardEntry{
			UserId:    entry.UserID,
			Rank:      entry.Rank,
			Xp:        entry.XP,
			LevelText: entry.LevelText,
		})
	}

	for _, reward := range snapshot.Rewards {
		response.Rewards = append(response.Rewards, &gamificationv1.RewardItem{
			Id:          reward.ID,
			Title:       reward.Title,
			Description: reward.Description,
			Cost:        reward.Cost,
			Status:      reward.Status,
			Category:    reward.Category,
		})
	}

	for _, purchase := range snapshot.Purchases {
		response.Purchases = append(response.Purchases, &gamificationv1.PurchaseItem{
			Id:         purchase.ID,
			RewardId:   purchase.RewardID,
			Title:      purchase.Title,
			Cost:       purchase.Cost,
			RedeemedAt: purchase.RedeemedAt,
			Status:     purchase.Status,
		})
	}

	for _, notification := range snapshot.Notifications {
		response.Notifications = append(response.Notifications, &gamificationv1.NotificationItem{
			Id:      notification.ID,
			Title:   notification.Title,
			Body:    notification.Body,
			Variant: notification.Variant,
		})
	}

	return response
}
