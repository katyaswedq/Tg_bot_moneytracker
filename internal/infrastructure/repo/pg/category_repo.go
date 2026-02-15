package pg

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	domainerr "tgfin/internal/domain/error"
	"tgfin/internal/domain/models"
	"tgfin/internal/domain/usecase"
)

type CategoryRepo struct{}

func NewCategoryRepo() *CategoryRepo {
	return &CategoryRepo{}
}

func (r *CategoryRepo) CreateDefault(ctx context.Context, exec user.Exec, userID uuid.UUID) error {
	query := `
		INSERT INTO categories (user_id, name, is_default)
		VALUES
			($1, 'Еда', TRUE),
			($1, 'Транспорт', TRUE),
			($1, 'Развлечения', TRUE),
			($1, 'Прочее', TRUE)
		ON CONFLICT (user_id, name) DO NOTHING
	`
	return exec.Exec(ctx, query, userID)
}

func (r *CategoryRepo) Create(ctx context.Context, exec user.Exec, userID uuid.UUID, name string) (*models.Category, error) {
	query := `
		INSERT INTO categories (user_id, name, is_default)
		VALUES ($1, $2, FALSE)
		RETURNING id, created_at, is_default, color
	`

	c := &models.Category{
		UserId:    userID,
		Name:      name,
		IsDefault: false,
	}

	err := exec.QueryRow(ctx, query, userID, name).Scan(&c.ID, &c.CreatedAt, &c.IsDefault, &c.Color)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, domainerr.ErrCategoryAlreadyExists
		}
		return nil, err
	}

	return c, nil
}

func (r *CategoryRepo) ListByUser(ctx context.Context, exec user.Exec, userID uuid.UUID) ([]*models.Category, error) {
	query := `
		SELECT id, user_id, name, color, created_at, is_default
		FROM categories
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at ASC
	`

	rows, err := exec.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*models.Category

	for rows.Next() {
		c := &models.Category{}
		if err := rows.Scan(&c.ID, &c.UserId, &c.Name, &c.Color, &c.CreatedAt, &c.IsDefault); err != nil {
			return nil, err
		}
		res = append(res, c)
	}

	return res, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, exec user.Exec, userID uuid.UUID, categoryID uuid.UUID) error {
	cmd, err := exec.ExecWithResult(ctx, `
		DELETE FROM categories
		WHERE id = $1 AND user_id = $2 AND is_default = FALSE
	`, categoryID, userID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domainerr.ErrCategoryNotFound
	}

	return nil
}

func (r *CategoryRepo) GetByName(ctx context.Context, exec user.Exec, userID uuid.UUID, name string) (*models.Category, error) {
	query := `
		SELECT id, user_id, name, color, created_at, is_default
		FROM categories
		WHERE user_id = $1 AND name = $2
		LIMIT 1
	`

	c := &models.Category{}
	err := exec.QueryRow(ctx, query, userID, name).Scan(&c.ID, &c.UserId, &c.Name, &c.Color, &c.CreatedAt, &c.IsDefault)
	if err != nil {
		return nil, domainerr.ErrCategoryNotFound
	}

	return c, nil
}
