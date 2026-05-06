package storage

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

func (r *Repository) List(ctx context.Context) ([]ColdStorage, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id::text, name, location_label, description, created_at, updated_at
		FROM cold_storages
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []ColdStorage{}
	for rows.Next() {
		item, err := scanColdStorage(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) Create(ctx context.Context, input CreateInput) (ColdStorage, error) {
	row := r.q.QueryRow(ctx, `
		INSERT INTO cold_storages (name, location_label, description)
		VALUES ($1, $2, $3)
		RETURNING id::text, name, location_label, description, created_at, updated_at
	`, input.Name, input.LocationLabel, input.Description)
	return scanColdStorage(row)
}

func scanColdStorage(row pgx.Row) (ColdStorage, error) {
	var item ColdStorage
	var location pgtype.Text
	var description pgtype.Text
	err := row.Scan(&item.ID, &item.Name, &location, &description, &item.CreatedAt, &item.UpdatedAt)
	item.LocationLabel = textPtr(location)
	item.Description = textPtr(description)
	return item, err
}

func textPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
