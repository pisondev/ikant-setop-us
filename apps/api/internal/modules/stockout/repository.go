package stockout

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type queryer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TxRepository interface {
	ListAvailableBatchesForUpdate(ctx context.Context, fishTypeID string) ([]AvailableBatch, error)
	CreateRecord(ctx context.Context, input CreateInput) (StockOut, error)
	CreateItem(ctx context.Context, stockOutID string, item Item) error
	UpdateBatchRemaining(ctx context.Context, id string, remainingWeightKG float64, status string) error
	CreateMovement(ctx context.Context, input MovementInput) error
}

type Store interface {
	TxRepository
	WithTx(ctx context.Context, fn func(TxRepository) error) error
	List(ctx context.Context, filter ListFilter) ([]StockOut, error)
}

type Repository struct {
	pool *pgxpool.Pool
	q    queryer
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, q: pool}
}

func (r *Repository) WithTx(ctx context.Context, fn func(TxRepository) error) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	txRepo := &Repository{q: tx}
	if err := fn(txRepo); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) ListAvailableBatchesForUpdate(ctx context.Context, fishTypeID string) ([]AvailableBatch, error) {
	rows, err := r.q.Query(ctx, `
		SELECT id::text, fish_type_id::text, remaining_weight_kg::float8, entered_at, status::text
		FROM stock_batches
		WHERE fish_type_id = $1
			AND status = 'available'
			AND remaining_weight_kg > 0
		ORDER BY entered_at ASC, created_at ASC
		FOR UPDATE
	`, fishTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []AvailableBatch{}
	for rows.Next() {
		var item AvailableBatch
		if err := rows.Scan(&item.ID, &item.FishTypeID, &item.RemainingWeightKG, &item.EnteredAt, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) CreateRecord(ctx context.Context, input CreateInput) (StockOut, error) {
	row := r.q.QueryRow(ctx, `
		INSERT INTO stock_outs (destination, total_weight_kg, out_at, notes)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, destination, total_weight_kg::float8, out_at, notes, created_at
	`, input.Destination, input.TotalWeightKG, input.OutAt, input.Notes)

	var item StockOut
	var notes pgtype.Text
	if err := row.Scan(&item.ID, &item.Destination, &item.TotalWeightKG, &item.OutAt, &notes, &item.CreatedAt); err != nil {
		return StockOut{}, err
	}
	item.Notes = textPtr(notes)
	return item, nil
}

func (r *Repository) CreateItem(ctx context.Context, stockOutID string, item Item) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO stock_out_items (stock_out_id, stock_batch_id, weight_kg)
		VALUES ($1, $2, $3)
	`, stockOutID, item.StockBatchID, item.WeightKG)
	return err
}

func (r *Repository) UpdateBatchRemaining(ctx context.Context, id string, remainingWeightKG float64, status string) error {
	_, err := r.q.Exec(ctx, `
		UPDATE stock_batches
		SET remaining_weight_kg = $2, status = $3, updated_at = now()
		WHERE id = $1
	`, id, remainingWeightKG, status)
	return err
}

func (r *Repository) CreateMovement(ctx context.Context, input MovementInput) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO stock_movements (stock_batch_id, movement_type, weight_kg, description)
		VALUES ($1, $2, $3, $4)
	`, input.StockBatchID, input.MovementType, input.WeightKG, input.Description)
	return err
}

func (r *Repository) List(ctx context.Context, filter ListFilter) ([]StockOut, error) {
	args := []any{}
	conditions := []string{"1 = 1"}

	addCondition := func(sql string, value any) {
		args = append(args, value)
		conditions = append(conditions, fmt.Sprintf(sql, len(args)))
	}
	if filter.FishTypeID != nil {
		addCondition(`EXISTS (
			SELECT 1
			FROM stock_out_items soi_filter
			JOIN stock_batches sb_filter ON sb_filter.id = soi_filter.stock_batch_id
			WHERE soi_filter.stock_out_id = so.id AND sb_filter.fish_type_id = $%d
		)`, *filter.FishTypeID)
	}
	if filter.Destination != nil {
		addCondition("so.destination ILIKE '%%' || $%d || '%%'", *filter.Destination)
	}
	if filter.DateFrom != nil {
		addCondition("so.out_at >= $%d", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		addCondition("so.out_at < $%d", filter.DateTo.Add(24*time.Hour))
	}

	query := fmt.Sprintf(`
		SELECT
			so.id::text,
			so.destination,
			so.total_weight_kg::float8,
			so.out_at,
			so.notes,
			so.created_at,
			soi.stock_batch_id::text,
			ft.name,
			soi.weight_kg::float8
		FROM stock_outs so
		LEFT JOIN stock_out_items soi ON soi.stock_out_id = so.id
		LEFT JOIN stock_batches sb ON sb.id = soi.stock_batch_id
		LEFT JOIN fish_types ft ON ft.id = sb.fish_type_id
		WHERE %s
		ORDER BY so.out_at DESC, so.created_at DESC, soi.created_at ASC
	`, strings.Join(conditions, " AND "))

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordered := []string{}
	byID := map[string]*StockOut{}
	for rows.Next() {
		var out StockOut
		var notes pgtype.Text
		var batchID pgtype.Text
		var fishTypeName pgtype.Text
		var weight *float64
		if err := rows.Scan(
			&out.ID,
			&out.Destination,
			&out.TotalWeightKG,
			&out.OutAt,
			&notes,
			&out.CreatedAt,
			&batchID,
			&fishTypeName,
			&weight,
		); err != nil {
			return nil, err
		}

		existing, ok := byID[out.ID]
		if !ok {
			out.Notes = textPtr(notes)
			out.Items = []Item{}
			byID[out.ID] = &out
			ordered = append(ordered, out.ID)
			existing = &out
		}
		if batchID.Valid && weight != nil {
			existing.Items = append(existing.Items, Item{
				StockBatchID: batchID.String,
				FishTypeName: textPtr(fishTypeName),
				WeightKG:     *weight,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	items := make([]StockOut, 0, len(ordered))
	for _, id := range ordered {
		items = append(items, *byID[id])
	}
	return items, nil
}

func textPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
