package grpcserver

import (
	"context"

	userv1 "cdek/platform/shared/contracts/user/v1"
	"cdek/platform/user-service/internal/application"
)

type Server struct {
	service *application.Service
}

func New(service *application.Service) *Server {
	return &Server{service: service}
}

func (s *Server) GetCurrentUser(_ context.Context, request *userv1.GetCurrentUserRequest) (*userv1.GetCurrentUserResponse, error) {
	user, err := s.service.GetCurrentUser(request.UserId)
	if err != nil {
		return nil, err
	}

	return &userv1.GetCurrentUserResponse{
		User: mapUser(user),
	}, nil
}

func (s *Server) BatchGetUsers(_ context.Context, request *userv1.BatchGetUsersRequest) (*userv1.BatchGetUsersResponse, error) {
	users, err := s.service.BatchGetUsers(request.UserIds)
	if err != nil {
		return nil, err
	}

	response := &userv1.BatchGetUsersResponse{
		Users: make([]*userv1.UserSummary, 0, len(users)),
	}

	for _, user := range users {
		response.Users = append(response.Users, mapUser(user))
	}

	return response, nil
}

func mapUser(user interface {
	GetID() string
	GetName() string
	GetTitle() string
	GetCompany() string
	GetLevel() int32
	GetLevelText() string
	GetJoinedAt() string
	GetLocation() string
	GetTeam() string
}) *userv1.UserSummary {
	return &userv1.UserSummary{
		Id:        user.GetID(),
		Name:      user.GetName(),
		Title:     user.GetTitle(),
		Company:   user.GetCompany(),
		Level:     user.GetLevel(),
		LevelText: user.GetLevelText(),
		JoinedAt:  user.GetJoinedAt(),
		Location:  user.GetLocation(),
		Team:      user.GetTeam(),
	}
}
