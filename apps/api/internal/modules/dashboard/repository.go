package dashboard

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func (r *Repository) Summary(ctx context.Context) (Summary, error) {
	var summary Summary
	err := r.q.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(remaining_weight_kg) FILTER (WHERE status = 'available'), 0)::float8,
			COUNT(*)::int,
			COUNT(*) FILTER (WHERE status = 'available')::int,
			COUNT(*) FILTER (WHERE status = 'depleted')::int
		FROM stock_batches
	`).Scan(
		&summary.TotalAvailableWeightKG,
		&summary.TotalStockBatches,
		&summary.TotalAvailableBatches,
		&summary.TotalDepletedBatches,
	)
	if err != nil {
		return Summary{}, err
	}

	if err := r.q.QueryRow(ctx, `
		SELECT COALESCE(SUM(weight_kg), 0)::float8
		FROM stock_movements
		WHERE movement_type = 'in' AND created_at >= CURRENT_DATE
	`).Scan(&summary.TodayStockInKG); err != nil {
		return Summary{}, err
	}

	if err := r.q.QueryRow(ctx, `
		SELECT COALESCE(SUM(total_weight_kg), 0)::float8
		FROM stock_outs
		WHERE created_at >= CURRENT_DATE
	`).Scan(&summary.TodayStockOutKG); err != nil {
		return Summary{}, err
	}

	fishRows, err := r.q.Query(ctx, `
		SELECT ft.id::text, ft.name, COALESCE(SUM(sb.remaining_weight_kg), 0)::float8, COUNT(sb.id)::int
		FROM fish_types ft
		JOIN stock_batches sb ON sb.fish_type_id = ft.id
			AND sb.status = 'available'
			AND sb.remaining_weight_kg > 0
		GROUP BY ft.id, ft.name
		ORDER BY ft.name ASC
	`)
	if err != nil {
		return Summary{}, err
	}
	defer fishRows.Close()
	summary.FishTypeSummary = []FishTypeSummary{}
	for fishRows.Next() {
		var item FishTypeSummary
		if err := fishRows.Scan(&item.FishTypeID, &item.FishTypeName, &item.AvailableWeightKG, &item.AvailableBatches); err != nil {
			return Summary{}, err
		}
		summary.FishTypeSummary = append(summary.FishTypeSummary, item)
	}
	if err := fishRows.Err(); err != nil {
		return Summary{}, err
	}

	storageRows, err := r.q.Query(ctx, `
		SELECT cs.id::text, cs.name, COALESCE(SUM(sb.remaining_weight_kg), 0)::float8, COUNT(sb.id)::int
		FROM cold_storages cs
		JOIN stock_batches sb ON sb.cold_storage_id = cs.id
			AND sb.status = 'available'
			AND sb.remaining_weight_kg > 0
		GROUP BY cs.id, cs.name
		ORDER BY cs.name ASC
	`)
	if err != nil {
		return Summary{}, err
	}
	defer storageRows.Close()
	summary.ColdStorageSummary = []ColdStorageSummary{}
	for storageRows.Next() {
		var item ColdStorageSummary
		if err := storageRows.Scan(&item.ColdStorageID, &item.ColdStorageName, &item.AvailableWeightKG, &item.AvailableBatches); err != nil {
			return Summary{}, err
		}
		summary.ColdStorageSummary = append(summary.ColdStorageSummary, item)
	}
	if err := storageRows.Err(); err != nil {
		return Summary{}, err
	}

	return summary, nil
}

func (r *Repository) RecentMovements(ctx context.Context, limit int) ([]RecentMovement, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	rows, err := r.q.Query(ctx, `
		SELECT
			sm.id::text,
			sm.stock_batch_id::text,
			sm.movement_type::text,
			ft.name,
			COALESCE(sm.weight_kg::float8, -1),
			COALESCE(sm.description, ''),
			sm.created_at
		FROM stock_movements sm
		JOIN stock_batches sb ON sb.id = sm.stock_batch_id
		JOIN fish_types ft ON ft.id = sb.fish_type_id
		ORDER BY sm.created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []RecentMovement{}
	for rows.Next() {
		var item RecentMovement
		var weight float64
		if err := rows.Scan(
			&item.ID,
			&item.StockBatchID,
			&item.MovementType,
			&item.FishTypeName,
			&weight,
			&item.Description,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		if weight >= 0 {
			item.WeightKG = &weight
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
