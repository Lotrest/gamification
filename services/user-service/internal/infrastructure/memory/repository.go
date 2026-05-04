package memory

import (
	"slices"

	"cdek/platform/user-service/internal/domain"
)

type Repository struct {
	users map[string]*domain.User
}

func NewRepository() *Repository {
	return &Repository{
		users: map[string]*domain.User{
			"me": {
				ID:        "me",
				Name:      "Алексей Воронов",
				Title:     "Frontend developer",
				Company:   "CDEK Digital",
				Level:     4,
				LevelText: "Lv.4",
				JoinedAt:  "Март 2023",
				Location:  "Новосибирск",
				Team:      "Платформа комьюнити",
			},
			"user-1": {
				ID:        "user-1",
				Name:      "Пользователь 1",
				Title:     "QA-инженер",
				Company:   "LogiTech Pro",
				Level:     18,
				LevelText: "Lv.18",
				JoinedAt:  "Январь 2021",
				Location:  "Москва",
				Team:      "QA Guild",
			},
			"user-2": {
				ID:        "user-2",
				Name:      "Пользователь 2",
				Title:     "Data Scientist",
				Company:   "DevStudio",
				Level:     16,
				LevelText: "Lv.16",
				JoinedAt:  "Июнь 2021",
				Location:  "Санкт-Петербург",
				Team:      "Analytics Lab",
			},
			"user-3": {
				ID:        "user-3",
				Name:      "Пользователь 3",
				Title:     "Web-разработчик",
				Company:   "TechSolutions LLC",
				Level:     12,
				LevelText: "Lv.12",
				JoinedAt:  "Август 2022",
				Location:  "Казань",
				Team:      "Frontend Platform",
			},
			"user-4": {
				ID:        "user-4",
				Name:      "Пользователь 4",
				Title:     "Backend developer",
				Company:   "CodeCraft",
				Level:     14,
				LevelText: "Lv.14",
				JoinedAt:  "Февраль 2022",
				Location:  "Екатеринбург",
				Team:      "Backend Core",
			},
		},
	}
}

func (r *Repository) GetCurrentUser(userID string) (*domain.User, error) {
	user, ok := r.users[userID]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *Repository) BatchGetUsers(userIDs []string) ([]*domain.User, error) {
	ordered := make([]*domain.User, 0, len(userIDs))

	for _, userID := range userIDs {
		if user, ok := r.users[userID]; ok {
			ordered = append(ordered, user)
		}
	}

	if len(ordered) == 0 && len(userIDs) > 0 {
		return nil, domain.ErrUserNotFound
	}

	slices.SortFunc(ordered, func(left, right *domain.User) int {
		switch {
		case left.Name < right.Name:
			return -1
		case left.Name > right.Name:
			return 1
		default:
			return 0
		}
	})

	return ordered, nil
}
