package stockout

import "time"

const (
	StatusAvailable = "available"
	StatusDepleted  = "depleted"
	MovementOut     = "out"
)

type AvailableBatch struct {
	ID                string
	FishTypeID        string
	RemainingWeightKG float64
	EnteredAt         time.Time
	Status            string
}

type Item struct {
	StockBatchID string  `json:"stock_batch_id"`
	FishTypeName *string `json:"fish_type_name,omitempty"`
	WeightKG     float64 `json:"weight_kg"`
}

type StockOut struct {
	ID            string    `json:"id"`
	FishTypeID    string    `json:"fish_type_id,omitempty"`
	Destination   string    `json:"destination"`
	TotalWeightKG float64   `json:"total_weight_kg"`
	OutAt         time.Time `json:"out_at"`
	Notes         *string   `json:"notes"`
	Items         []Item    `json:"items"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateInput struct {
	FishTypeID    string
	TotalWeightKG float64
	Destination   string
	OutAt         time.Time
	Notes         *string
}

type MovementInput struct {
	StockBatchID string
	MovementType string
	WeightKG     *float64
	Description  string
}

type ListFilter struct {
	FishTypeID  *string
	Destination *string
	DateFrom    *time.Time
	DateTo      *time.Time
}
