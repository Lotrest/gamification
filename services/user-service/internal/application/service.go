package application

import "cdek/platform/user-service/internal/domain"

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetCurrentUser(userID string) (*domain.User, error) {
	return s.repository.GetCurrentUser(userID)
}

func (s *Service) BatchGetUsers(userIDs []string) ([]*domain.User, error) {
	return s.repository.BatchGetUsers(userIDs)
}
