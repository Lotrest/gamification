package v1

import (
	"context"

	"google.golang.org/grpc"
)

const (
	GamificationServiceName      = "cdek.gamification.v1.GamificationService"
	GetPortalSnapshotMethodName  = "/cdek.gamification.v1.GamificationService/GetPortalSnapshot"
	AcceptTaskMethodName         = "/cdek.gamification.v1.GamificationService/AcceptTask"
	AdvanceTaskMethodName        = "/cdek.gamification.v1.GamificationService/AdvanceTask"
	RedeemRewardMethodName       = "/cdek.gamification.v1.GamificationService/RedeemReward"
)

type GetPortalSnapshotRequest struct {
	UserId string `json:"userId"`
}

type AcceptTaskRequest struct {
	UserId string `json:"userId"`
	TaskId string `json:"taskId"`
}

type AdvanceTaskRequest struct {
	UserId string `json:"userId"`
	TaskId string `json:"taskId"`
}

type RedeemRewardRequest struct {
	UserId   string `json:"userId"`
	RewardId string `json:"rewardId"`
}

type MutationResponse struct {
	Snapshot *GetPortalSnapshotResponse `json:"snapshot"`
}

type UserSummary struct {
	UserId           string `json:"userId"`
	CurrentXp        int32  `json:"currentXp"`
	Level            int32  `json:"level"`
	LevelText        string `json:"levelText"`
	ProgressPercent  int32  `json:"progressPercent"`
	XpToNextLevel    int32  `json:"xpToNextLevel"`
	Coins            int32  `json:"coins"`
	TodayEarned      int32  `json:"todayEarned"`
	StreakDays       int32  `json:"streakDays"`
	Rank             int32  `json:"rank"`
	TotalBadges      int32  `json:"totalBadges"`
	CompletedTasks   int32  `json:"completedTasks"`
}

type MetricCard struct {
	Id      string `json:"id"`
	Label   string `json:"label"`
	Value   string `json:"value"`
	Caption string `json:"caption"`
}

type ActivityPoint struct {
	Day string `json:"day"`
	Xp  int32  `json:"xp"`
}

type ArticleCard struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Views    int32  `json:"views"`
	Comments int32  `json:"comments"`
	Xp       int32  `json:"xp"`
	Rating   string `json:"rating"`
}

type ActivityItem struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"`
	Xp        int32  `json:"xp"`
}

type TaskItem struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Progress    int32  `json:"progress"`
	Target      int32  `json:"target"`
	RewardXp    int32  `json:"rewardXp"`
}

type AchievementBucket struct {
	Id         string `json:"id"`
	Label      string `json:"label"`
	Collected  int32  `json:"collected"`
	Total      int32  `json:"total"`
	Accent     string `json:"accent"`
}

type AchievementItem struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Rarity      string `json:"rarity"`
	Status      string `json:"status"`
	RewardXp    int32  `json:"rewardXp"`
}

type LeaderboardEntry struct {
	UserId    string `json:"userId"`
	Rank      int32  `json:"rank"`
	Xp        int32  `json:"xp"`
	LevelText string `json:"levelText"`
}

type RewardItem struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Cost        int32  `json:"cost"`
	Status      string `json:"status"`
	Category    string `json:"category"`
}

type PurchaseItem struct {
	Id         string `json:"id"`
	RewardId   string `json:"rewardId"`
	Title      string `json:"title"`
	Cost       int32  `json:"cost"`
	RedeemedAt string `json:"redeemedAt"`
	Status     string `json:"status"`
}

type NotificationItem struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Variant string `json:"variant"`
}

type GetPortalSnapshotResponse struct {
	Summary            *UserSummary         `json:"summary"`
	Metrics            []*MetricCard        `json:"metrics"`
	WeeklyActivity     []*ActivityPoint     `json:"weeklyActivity"`
	Articles           []*ArticleCard       `json:"articles"`
	RecentActivity     []*ActivityItem      `json:"recentActivity"`
	Tasks              []*TaskItem          `json:"tasks"`
	AchievementBuckets []*AchievementBucket `json:"achievementBuckets"`
	Achievements       []*AchievementItem   `json:"achievements"`
	Leaderboard        []*LeaderboardEntry  `json:"leaderboard"`
	Rewards            []*RewardItem        `json:"rewards"`
	Purchases          []*PurchaseItem      `json:"purchases"`
	Notifications      []*NotificationItem  `json:"notifications"`
}

type GamificationServiceClient interface {
	GetPortalSnapshot(ctx context.Context, in *GetPortalSnapshotRequest, opts ...grpc.CallOption) (*GetPortalSnapshotResponse, error)
	AcceptTask(ctx context.Context, in *AcceptTaskRequest, opts ...grpc.CallOption) (*MutationResponse, error)
	AdvanceTask(ctx context.Context, in *AdvanceTaskRequest, opts ...grpc.CallOption) (*MutationResponse, error)
	RedeemReward(ctx context.Context, in *RedeemRewardRequest, opts ...grpc.CallOption) (*MutationResponse, error)
}

type gamificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGamificationServiceClient(cc grpc.ClientConnInterface) GamificationServiceClient {
	return &gamificationServiceClient{cc: cc}
}

func (c *gamificationServiceClient) GetPortalSnapshot(ctx context.Context, in *GetPortalSnapshotRequest, opts ...grpc.CallOption) (*GetPortalSnapshotResponse, error) {
	out := new(GetPortalSnapshotResponse)
	if err := c.cc.Invoke(ctx, GetPortalSnapshotMethodName, in, out, opts...); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *gamificationServiceClient) AcceptTask(ctx context.Context, in *AcceptTaskRequest, opts ...grpc.CallOption) (*MutationResponse, error) {
	out := new(MutationResponse)
	if err := c.cc.Invoke(ctx, AcceptTaskMethodName, in, out, opts...); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *gamificationServiceClient) AdvanceTask(ctx context.Context, in *AdvanceTaskRequest, opts ...grpc.CallOption) (*MutationResponse, error) {
	out := new(MutationResponse)
	if err := c.cc.Invoke(ctx, AdvanceTaskMethodName, in, out, opts...); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *gamificationServiceClient) RedeemReward(ctx context.Context, in *RedeemRewardRequest, opts ...grpc.CallOption) (*MutationResponse, error) {
	out := new(MutationResponse)
	if err := c.cc.Invoke(ctx, RedeemRewardMethodName, in, out, opts...); err != nil {
		return nil, err
	}

	return out, nil
}

type GamificationServiceServer interface {
	GetPortalSnapshot(context.Context, *GetPortalSnapshotRequest) (*GetPortalSnapshotResponse, error)
	AcceptTask(context.Context, *AcceptTaskRequest) (*MutationResponse, error)
	AdvanceTask(context.Context, *AdvanceTaskRequest) (*MutationResponse, error)
	RedeemReward(context.Context, *RedeemRewardRequest) (*MutationResponse, error)
}

func RegisterGamificationServiceServer(s grpc.ServiceRegistrar, srv GamificationServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: GamificationServiceName,
		HandlerType: (*GamificationServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetPortalSnapshot",
				Handler:    getPortalSnapshotHandler(srv),
			},
			{
				MethodName: "AcceptTask",
				Handler:    acceptTaskHandler(srv),
			},
			{
				MethodName: "AdvanceTask",
				Handler:    advanceTaskHandler(srv),
			},
			{
				MethodName: "RedeemReward",
				Handler:    redeemRewardHandler(srv),
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "proto/gamification/v1/gamification.proto",
	}, srv)
}

func getPortalSnapshotHandler(srv GamificationServiceServer) grpc.MethodHandler {
	return func(service any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
		in := new(GetPortalSnapshotRequest)
		if err := dec(in); err != nil {
			return nil, err
		}

		if interceptor == nil {
			return srv.GetPortalSnapshot(ctx, in)
		}

		info := &grpc.UnaryServerInfo{
			Server:     service,
			FullMethod: GetPortalSnapshotMethodName,
		}

		handler := func(ctx context.Context, req any) (any, error) {
			return srv.GetPortalSnapshot(ctx, req.(*GetPortalSnapshotRequest))
		}

		return interceptor(ctx, in, info, handler)
	}
}

func acceptTaskHandler(srv GamificationServiceServer) grpc.MethodHandler {
	return func(service any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
		in := new(AcceptTaskRequest)
		if err := dec(in); err != nil {
			return nil, err
		}

		if interceptor == nil {
			return srv.AcceptTask(ctx, in)
		}

		info := &grpc.UnaryServerInfo{
			Server:     service,
			FullMethod: AcceptTaskMethodName,
		}

		handler := func(ctx context.Context, req any) (any, error) {
			return srv.AcceptTask(ctx, req.(*AcceptTaskRequest))
		}

		return interceptor(ctx, in, info, handler)
	}
}

func advanceTaskHandler(srv GamificationServiceServer) grpc.MethodHandler {
	return func(service any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
		in := new(AdvanceTaskRequest)
		if err := dec(in); err != nil {
			return nil, err
		}

		if interceptor == nil {
			return srv.AdvanceTask(ctx, in)
		}

		info := &grpc.UnaryServerInfo{
			Server:     service,
			FullMethod: AdvanceTaskMethodName,
		}

		handler := func(ctx context.Context, req any) (any, error) {
			return srv.AdvanceTask(ctx, req.(*AdvanceTaskRequest))
		}

		return interceptor(ctx, in, info, handler)
	}
}

func redeemRewardHandler(srv GamificationServiceServer) grpc.MethodHandler {
	return func(service any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
		in := new(RedeemRewardRequest)
		if err := dec(in); err != nil {
			return nil, err
		}

		if interceptor == nil {
			return srv.RedeemReward(ctx, in)
		}

		info := &grpc.UnaryServerInfo{
			Server:     service,
			FullMethod: RedeemRewardMethodName,
		}

		handler := func(ctx context.Context, req any) (any, error) {
			return srv.RedeemReward(ctx, req.(*RedeemRewardRequest))
		}

		return interceptor(ctx, in, info, handler)
	}
}
