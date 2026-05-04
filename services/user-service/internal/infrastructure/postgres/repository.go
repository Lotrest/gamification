package postgres

import (
	"context"
	"errors"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"cdek/platform/user-service/internal/domain"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetCurrentUser(userID string) (*domain.User, error) {
	const query = `
		select id, name, title, company, level, level_text, joined_at, location, team
		from user_service.users
		where id = $1
	`

	var user domain.User
	err := r.pool.QueryRow(context.Background(), query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Title,
		&user.Company,
		&user.Level,
		&user.LevelText,
		&user.JoinedAt,
		&user.Location,
		&user.Team,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) BatchGetUsers(userIDs []string) ([]*domain.User, error) {
	const query = `
		select id, name, title, company, level, level_text, joined_at, location, team
		from user_service.users
		where id = any($1::text[])
	`

	rows, err := r.pool.Query(context.Background(), query, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0, len(userIDs))
	for rows.Next() {
		var user domain.User
		if scanErr := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Title,
			&user.Company,
			&user.Level,
			&user.LevelText,
			&user.JoinedAt,
			&user.Location,
			&user.Team,
		); scanErr != nil {
			return nil, scanErr
		}

		users = append(users, &user)
	}

	if len(users) == 0 && len(userIDs) > 0 {
		return nil, domain.ErrUserNotFound
	}

	sort.Slice(users, func(left, right int) bool {
		return users[left].Name < users[right].Name
	})

	return users, rows.Err()
}
