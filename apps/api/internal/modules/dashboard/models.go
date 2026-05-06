package dashboard

import "time"

type FishTypeSummary struct {
	FishTypeID        string  `json:"fish_type_id"`
	FishTypeName      string  `json:"fish_type_name"`
	AvailableWeightKG float64 `json:"available_weight_kg"`
	AvailableBatches  int     `json:"available_batches"`
}

type ColdStorageSummary struct {
	ColdStorageID     string  `json:"cold_storage_id"`
	ColdStorageName   string  `json:"cold_storage_name"`
	AvailableWeightKG float64 `json:"available_weight_kg"`
	AvailableBatches  int     `json:"available_batches"`
}

type Summary struct {
	TotalAvailableWeightKG float64              `json:"total_available_weight_kg"`
	TotalStockBatches      int                  `json:"total_stock_batches"`
	TotalAvailableBatches  int                  `json:"total_available_batches"`
	TotalDepletedBatches   int                  `json:"total_depleted_batches"`
	TodayStockInKG         float64              `json:"today_stock_in_kg"`
	TodayStockOutKG        float64              `json:"today_stock_out_kg"`
	FishTypeSummary        []FishTypeSummary    `json:"fish_type_summary"`
	ColdStorageSummary     []ColdStorageSummary `json:"cold_storage_summary"`
}

type RecentMovement struct {
	ID           string    `json:"id"`
	StockBatchID string    `json:"stock_batch_id"`
	MovementType string    `json:"movement_type"`
	FishTypeName string    `json:"fish_type_name"`
	WeightKG     *float64  `json:"weight_kg"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}
