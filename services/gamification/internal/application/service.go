package application

import (
	"fmt"

	"cdek/platform/gamification/internal/domain"
)

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetSnapshot(userID string) (*domain.Snapshot, error) {
	state, err := s.repository.GetState(userID)
	if err != nil {
		return nil, err
	}

	return buildSnapshot(state), nil
}

func (s *Service) AcceptTask(userID, taskID string) (*domain.Snapshot, error) {
	state, err := s.repository.AcceptTask(userID, taskID)
	if err != nil {
		return nil, err
	}

	return buildSnapshot(state), nil
}

func (s *Service) AdvanceTask(userID, taskID string) (*domain.Snapshot, error) {
	state, err := s.repository.AdvanceTask(userID, taskID)
	if err != nil {
		return nil, err
	}

	return buildSnapshot(state), nil
}

func (s *Service) RedeemReward(userID, rewardID string) (*domain.Snapshot, error) {
	state, err := s.repository.RedeemReward(userID, rewardID)
	if err != nil {
		return nil, err
	}

	return buildSnapshot(state), nil
}

func buildSnapshot(state *domain.PortalState) *domain.Snapshot {
	totalBadges := int32(0)
	for _, achievement := range state.Achievements {
		if achievement.Status == domain.AchievementStatusUnlocked {
			totalBadges++
		}
	}

	level, progressPercent, xpToNext := computeLevel(state.CurrentXP)
	summary := domain.Summary{
		UserID:          state.UserID,
		CurrentXP:       state.CurrentXP,
		Level:           level,
		LevelText:       fmt.Sprintf("Lv.%d", level),
		ProgressPercent: progressPercent,
		XpToNextLevel:   xpToNext,
		Coins:           state.Coins,
		TodayEarned:     state.TodayEarned,
		StreakDays:      state.StreakDays,
		Rank:            state.Rank,
		TotalBadges:     totalBadges,
		CompletedTasks:  state.CompletedTasks,
	}

	return &domain.Snapshot{
		Summary: summary,
		Metrics: []domain.MetricCard{
			{ID: "api", Label: "API запросы", Value: fmt.Sprintf("%d", state.ApiRequests), Caption: "за 30 дней"},
			{ID: "articles", Label: "Статьи", Value: fmt.Sprintf("%d", state.ArticlesCount), Caption: "опубликовано"},
			{ID: "comments", Label: "Комментарии", Value: fmt.Sprintf("%d", state.CommentsCount), Caption: "в обсуждениях"},
			{ID: "xp", Label: "Всего XP", Value: fmt.Sprintf("%d", state.CurrentXP), Caption: "накопленный опыт"},
		},
		WeeklyActivity:     append([]domain.ActivityPoint(nil), state.WeeklyActivity...),
		Articles:           append([]domain.ArticleCard(nil), state.Articles...),
		RecentActivity:     append([]domain.ActivityItem(nil), state.RecentActivity...),
		Tasks:              append([]domain.Task(nil), state.Tasks...),
		AchievementBuckets: buildAchievementBuckets(state.Achievements),
		Achievements:       append([]domain.Achievement(nil), state.Achievements...),
		Leaderboard:        enrichLeaderboard(state.Leaderboard),
		Rewards:            append([]domain.Reward(nil), state.Rewards...),
		Purchases:          append([]domain.Purchase(nil), state.Purchases...),
		Notifications:      append([]domain.Notification(nil), state.Notifications...),
	}
}

func computeLevel(currentXP int32) (int32, int32, int32) {
	thresholds := []int32{0, 500, 1200, 2200, 3600, 5400, 7600, 10200, 13200, 16800}

	level := int32(1)
	for index := 1; index < len(thresholds); index++ {
		if currentXP >= thresholds[index] {
			level = int32(index + 1)
		}
	}

	currentIndex := level - 1
	if int(currentIndex) >= len(thresholds)-1 {
		return level, 100, 0
	}

	currentBase := thresholds[currentIndex]
	nextBase := thresholds[currentIndex+1]
	span := nextBase - currentBase
	progress := currentXP - currentBase

	return level, (progress * 100) / span, nextBase - currentXP
}

func buildAchievementBuckets(achievements []domain.Achievement) []domain.AchievementBucket {
	type bucket struct {
		label string
		total int32
		open  int32
		color string
	}

	buckets := map[string]*bucket{
		"Обычные":      {label: "Обычные", color: "green"},
		"Редкие":       {label: "Редкие", color: "gray"},
		"Эпические":    {label: "Эпические", color: "green"},
		"Легендарные":  {label: "Легендарные", color: "gray"},
	}

	for _, achievement := range achievements {
		current, ok := buckets[achievement.Rarity]
		if !ok {
			continue
		}

		current.total++
		if achievement.Status == domain.AchievementStatusUnlocked {
			current.open++
		}
	}

	return []domain.AchievementBucket{
		{ID: "common", Label: buckets["Обычные"].label, Collected: buckets["Обычные"].open, Total: buckets["Обычные"].total, Accent: buckets["Обычные"].color},
		{ID: "rare", Label: buckets["Редкие"].label, Collected: buckets["Редкие"].open, Total: buckets["Редкие"].total, Accent: buckets["Редкие"].color},
		{ID: "epic", Label: buckets["Эпические"].label, Collected: buckets["Эпические"].open, Total: buckets["Эпические"].total, Accent: buckets["Эпические"].color},
		{ID: "legendary", Label: buckets["Легендарные"].label, Collected: buckets["Легендарные"].open, Total: buckets["Легендарные"].total, Accent: buckets["Легендарные"].color},
	}
}

func enrichLeaderboard(entries []domain.LeaderboardEntry) []domain.LeaderboardEntry {
	enriched := make([]domain.LeaderboardEntry, 0, len(entries))
	for _, entry := range entries {
		level, _, _ := computeLevel(entry.XP)
		entry.LevelText = fmt.Sprintf("Lv.%d", level)
		enriched = append(enriched, entry)
	}

	return enriched
}
