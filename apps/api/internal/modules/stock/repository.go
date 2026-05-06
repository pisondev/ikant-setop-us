package stock

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("stock batch not found")

type queryer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TxRepository interface {
	CreateBatch(ctx context.Context, input CreateInput) (Batch, error)
	CreateMovement(ctx context.Context, input MovementInput) error
	GetBatch(ctx context.Context, id string) (Batch, error)
	UpdateBatchLocation(ctx context.Context, id string, coldStorageID string) error
	UpdateBatchQuality(ctx context.Context, id string, quality string) error
}

type Store interface {
	TxRepository
	WithTx(ctx context.Context, fn func(TxRepository) error) error
	List(ctx context.Context, filter ListFilter) ([]Detail, error)
	GetDetail(ctx context.Context, id string) (Detail, error)
	ListFIFO(ctx context.Context, filter FIFOFilter) ([]FIFOItem, error)
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

func (r *Repository) List(ctx context.Context, filter ListFilter) ([]Detail, error) {
	args := []any{}
	conditions := []string{"1 = 1"}

	addCondition := func(sql string, value any) {
		args = append(args, value)
		conditions = append(conditions, fmt.Sprintf(sql, len(args)))
	}
	if filter.FishTypeID != nil {
		addCondition("sb.fish_type_id = $%d", *filter.FishTypeID)
	}
	if filter.Quality != nil {
		addCondition("sb.quality = $%d", *filter.Quality)
	}
	if filter.ColdStorageID != nil {
		addCondition("sb.cold_storage_id = $%d", *filter.ColdStorageID)
	}
	if filter.Status != nil {
		addCondition("sb.status = $%d", *filter.Status)
	}

	orderBy := "sb.created_at DESC"
	if filter.Sort == "" || filter.Sort == "fifo" {
		orderBy = "sb.entered_at ASC, sb.created_at ASC"
	}

	query := fmt.Sprintf(`
		SELECT
			sb.id::text,
			ft.id::text,
			ft.name,
			cs.id::text,
			cs.name,
			cs.location_label,
			sb.quality::text,
			sb.initial_weight_kg::float8,
			sb.remaining_weight_kg::float8,
			sb.entered_at,
			sb.status::text,
			sb.notes,
			sb.created_at,
			sb.updated_at
		FROM stock_batches sb
		JOIN fish_types ft ON ft.id = sb.fish_type_id
		JOIN cold_storages cs ON cs.id = sb.cold_storage_id
		WHERE %s
		ORDER BY %s
	`, strings.Join(conditions, " AND "), orderBy)

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Detail{}
	for rows.Next() {
		item, err := scanDetail(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) CreateBatch(ctx context.Context, input CreateInput) (Batch, error) {
	row := r.q.QueryRow(ctx, `
		INSERT INTO stock_batches (
			fish_type_id,
			cold_storage_id,
			quality,
			initial_weight_kg,
			remaining_weight_kg,
			entered_at,
			status,
			notes
		)
		VALUES ($1, $2, $3, $4, $4, $5, 'available', $6)
		RETURNING id::text, fish_type_id::text, cold_storage_id::text, quality::text,
			initial_weight_kg::float8, remaining_weight_kg::float8, entered_at,
			status::text, notes, created_at, updated_at
	`, input.FishTypeID, input.ColdStorageID, input.Quality, input.InitialWeightKG, input.EnteredAt, input.Notes)
	return scanBatch(row)
}

func (r *Repository) GetDetail(ctx context.Context, id string) (Detail, error) {
	row := r.q.QueryRow(ctx, `
		SELECT
			sb.id::text,
			ft.id::text,
			ft.name,
			cs.id::text,
			cs.name,
			cs.location_label,
			sb.quality::text,
			sb.initial_weight_kg::float8,
			sb.remaining_weight_kg::float8,
			sb.entered_at,
			sb.status::text,
			sb.notes,
			sb.created_at,
			sb.updated_at
		FROM stock_batches sb
		JOIN fish_types ft ON ft.id = sb.fish_type_id
		JOIN cold_storages cs ON cs.id = sb.cold_storage_id
		WHERE sb.id = $1
	`, id)
	item, err := scanDetail(row)
	return item, mapNoRows(err)
}

func (r *Repository) GetBatch(ctx context.Context, id string) (Batch, error) {
	row := r.q.QueryRow(ctx, `
		SELECT id::text, fish_type_id::text, cold_storage_id::text, quality::text,
			initial_weight_kg::float8, remaining_weight_kg::float8, entered_at,
			status::text, notes, created_at, updated_at
		FROM stock_batches
		WHERE id = $1
	`, id)
	item, err := scanBatch(row)
	return item, mapNoRows(err)
}

func (r *Repository) UpdateBatchQuality(ctx context.Context, id string, quality string) error {
	tag, err := r.q.Exec(ctx, `
		UPDATE stock_batches
		SET quality = $2, updated_at = now()
		WHERE id = $1
	`, id, quality)
	return mapUpdateResult(tag, err)
}

func (r *Repository) UpdateBatchLocation(ctx context.Context, id string, coldStorageID string) error {
	tag, err := r.q.Exec(ctx, `
		UPDATE stock_batches
		SET cold_storage_id = $2, updated_at = now()
		WHERE id = $1
	`, id, coldStorageID)
	return mapUpdateResult(tag, err)
}

func (r *Repository) ListFIFO(ctx context.Context, filter FIFOFilter) ([]FIFOItem, error) {
	args := []any{}
	conditions := []string{"sb.status = 'available'", "sb.remaining_weight_kg > 0"}
	if filter.FishTypeID != nil {
		args = append(args, *filter.FishTypeID)
		conditions = append(conditions, fmt.Sprintf("sb.fish_type_id = $%d", len(args)))
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	args = append(args, limit, offset)
	limitArg := len(args) - 1
	offsetArg := len(args)

	query := fmt.Sprintf(`
		SELECT
			sb.id::text,
			ft.name,
			sb.quality::text,
			sb.remaining_weight_kg::float8,
			sb.entered_at,
			cs.name,
			cs.location_label,
			ROW_NUMBER() OVER (ORDER BY sb.entered_at ASC, sb.created_at ASC)::int AS fifo_rank
		FROM stock_batches sb
		JOIN fish_types ft ON ft.id = sb.fish_type_id
		JOIN cold_storages cs ON cs.id = sb.cold_storage_id
		WHERE %s
		ORDER BY sb.entered_at ASC, sb.created_at ASC
		LIMIT $%d OFFSET $%d
	`, strings.Join(conditions, " AND "), limitArg, offsetArg)

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []FIFOItem{}
	for rows.Next() {
		var item FIFOItem
		var location pgtype.Text
		if err := rows.Scan(
			&item.ID,
			&item.FishTypeName,
			&item.Quality,
			&item.RemainingWeightKG,
			&item.EnteredAt,
			&item.ColdStorageName,
			&location,
			&item.FIFORank,
		); err != nil {
			return nil, err
		}
		item.LocationLabel = textPtr(location)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) CreateMovement(ctx context.Context, input MovementInput) error {
	_, err := r.q.Exec(ctx, `
		INSERT INTO stock_movements (
			stock_batch_id,
			movement_type,
			weight_kg,
			previous_quality,
			new_quality,
			previous_cold_storage_id,
			new_cold_storage_id,
			description
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, input.StockBatchID, input.MovementType, input.WeightKG, input.PreviousQuality, input.NewQuality, input.PreviousColdStorageID, input.NewColdStorageID, input.Description)
	return err
}

func scanBatch(row pgx.Row) (Batch, error) {
	var item Batch
	var notes pgtype.Text
	err := row.Scan(
		&item.ID,
		&item.FishTypeID,
		&item.ColdStorageID,
		&item.Quality,
		&item.InitialWeightKG,
		&item.RemainingWeightKG,
		&item.EnteredAt,
		&item.Status,
		&notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	item.Notes = textPtr(notes)
	return item, mapNoRows(err)
}

func scanDetail(row pgx.Row) (Detail, error) {
	var item Detail
	var location pgtype.Text
	var notes pgtype.Text
	err := row.Scan(
		&item.ID,
		&item.FishType.ID,
		&item.FishType.Name,
		&item.ColdStorage.ID,
		&item.ColdStorage.Name,
		&location,
		&item.Quality,
		&item.InitialWeightKG,
		&item.RemainingWeightKG,
		&item.EnteredAt,
		&item.Status,
		&notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	item.ColdStorage.LocationLabel = textPtr(location)
	item.Notes = textPtr(notes)
	return item, mapNoRows(err)
}

func mapNoRows(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func mapUpdateResult(tag pgconn.CommandTag, err error) error {
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func textPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
