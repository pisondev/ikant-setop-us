package fish

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository struct {
	q queryer
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{q: pool}
}

func (r *Repository) List(ctx context.Context) ([]FishType, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id::text, name, image_url, description, created_at, updated_at
		FROM fish_types
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []FishType{}
	for rows.Next() {
		item, err := scanFishType(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) Create(ctx context.Context, input CreateInput) (FishType, error) {
	row := r.q.QueryRow(ctx, `
		INSERT INTO fish_types (name, image_url, description)
		VALUES ($1, $2, $3)
		RETURNING id::text, name, image_url, description, created_at, updated_at
	`, input.Name, input.ImageURL, input.Description)
	return scanFishType(row)
}

func scanFishType(row pgx.Row) (FishType, error) {
	var item FishType
	var imageURL pgtype.Text
	var description pgtype.Text
	err := row.Scan(&item.ID, &item.Name, &imageURL, &description, &item.CreatedAt, &item.UpdatedAt)
	item.ImageURL = textPtr(imageURL)
	item.Description = textPtr(description)
	return item, err
}

func textPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
