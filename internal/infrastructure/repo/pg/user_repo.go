package pg

import (
	"context"

	"tgfin/internal/domain/models"
	domainErr "tgfin/internal/domain/error"
	"tgfin/internal/domain/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo { 
	return &UserRepo{} 
}

func (r *UserRepo) Create(ctx context.Context, exec user.Exec, u *models.User) error {
	query := `
		INSERT INTO users (telegram_id, user_name, first_name, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := exec.QueryRow(ctx, query, u.TelegramID, u.UserName, u.FirstName, u.CreatedAt).Scan(&u.ID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return domainErr.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserRepo) GetByTelegramID(ctx context.Context, exec user.Exec, telegramID int64) (*models.User, error) {
	query := `
		SELECT id, telegram_id, user_name, first_name, created_at
		FROM users
		WHERE telegram_id = $1
	`

	u := &models.User{}
	err := exec.QueryRow(ctx, query, telegramID).Scan(
		&u.ID,
		&u.TelegramID,
		&u.UserName,
		&u.FirstName,
		&u.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domainErr.ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}
