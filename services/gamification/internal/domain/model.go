package domain

import "errors"

var (
	ErrTaskNotFound        = errors.New("task not found")
	ErrRewardNotFound      = errors.New("reward not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

const (
	TaskStatusAvailable  = "available"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"

	AchievementStatusUnlocked = "unlocked"
	AchievementStatusLocked   = "locked"

	RewardStatusAvailable = "available"
	RewardStatusRedeemed  = "redeemed"
)

type Summary struct {
	UserID          string
	CurrentXP       int32
	Level           int32
	LevelText       string
	ProgressPercent int32
	XpToNextLevel   int32
	Coins           int32
	TodayEarned     int32
	StreakDays      int32
	Rank            int32
	TotalBadges     int32
	CompletedTasks  int32
}

type MetricCard struct {
	ID      string
	Label   string
	Value   string
	Caption string
}

type ActivityPoint struct {
	Day string
	XP  int32
}

type ArticleCard struct {
	ID       string
	Title    string
	Views    int32
	Comments int32
	XP       int32
	Rating   string
}

type ActivityItem struct {
	ID        string
	Title     string
	Timestamp string
	XP        int32
}

type Task struct {
	ID          string
	Title       string
	Description string
	Status      string
	Progress    int32
	Target      int32
	RewardXP    int32
}

type AchievementBucket struct {
	ID        string
	Label     string
	Collected int32
	Total     int32
	Accent    string
}

type Achievement struct {
	ID          string
	Title       string
	Description string
	Rarity      string
	Status      string
	RewardXP    int32
}

type LeaderboardEntry struct {
	UserID    string
	Rank      int32
	XP        int32
	LevelText string
}

type Reward struct {
	ID          string
	Title       string
	Description string
	Cost        int32
	Status      string
	Category    string
}

type Purchase struct {
	ID         string
	RewardID   string
	Title      string
	Cost       int32
	RedeemedAt string
	Status     string
}

type Notification struct {
	ID      string
	Title   string
	Body    string
	Variant string
}

type PortalState struct {
	UserID          string
	CurrentXP       int32
	Coins           int32
	TodayEarned     int32
	StreakDays      int32
	Rank            int32
	CompletedTasks  int32
	ApiRequests     int32
	ArticlesCount   int32
	CommentsCount   int32
	WeeklyActivity  []ActivityPoint
	Articles        []ArticleCard
	RecentActivity  []ActivityItem
	Tasks           []Task
	Achievements    []Achievement
	Leaderboard     []LeaderboardEntry
	Rewards         []Reward
	Purchases       []Purchase
	Notifications   []Notification
}

type Snapshot struct {
	Summary            Summary
	Metrics            []MetricCard
	WeeklyActivity     []ActivityPoint
	Articles           []ArticleCard
	RecentActivity     []ActivityItem
	Tasks              []Task
	AchievementBuckets []AchievementBucket
	Achievements       []Achievement
	Leaderboard        []LeaderboardEntry
	Rewards            []Reward
	Purchases          []Purchase
	Notifications      []Notification
}

type Repository interface {
	GetState(userID string) (*PortalState, error)
	AcceptTask(userID, taskID string) (*PortalState, error)
	AdvanceTask(userID, taskID string) (*PortalState, error)
	RedeemReward(userID, rewardID string) (*PortalState, error)
}
