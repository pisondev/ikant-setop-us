package stock

import "time"

const (
	QualityGood   = "baik"
	QualityMedium = "sedang"
	QualityBad    = "buruk"

	StatusAvailable = "available"
	StatusDepleted  = "depleted"

	MovementIn             = "in"
	MovementQualityUpdate  = "quality_update"
	MovementLocationUpdate = "location_update"
)

type Batch struct {
	ID                string    `json:"id"`
	FishTypeID        string    `json:"fish_type_id"`
	ColdStorageID     string    `json:"cold_storage_id"`
	Quality           string    `json:"quality"`
	InitialWeightKG   float64   `json:"initial_weight_kg"`
	RemainingWeightKG float64   `json:"remaining_weight_kg"`
	EnteredAt         time.Time `json:"entered_at"`
	Status            string    `json:"status"`
	Notes             *string   `json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type FishType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ColdStorage struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	LocationLabel *string `json:"location_label"`
}

type Detail struct {
	ID                string      `json:"id"`
	FishType          FishType    `json:"fish_type"`
	ColdStorage       ColdStorage `json:"cold_storage"`
	Quality           string      `json:"quality"`
	InitialWeightKG   float64     `json:"initial_weight_kg"`
	RemainingWeightKG float64     `json:"remaining_weight_kg"`
	EnteredAt         time.Time   `json:"entered_at"`
	Status            string      `json:"status"`
	Notes             *string     `json:"notes"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type FIFOItem struct {
	ID                string    `json:"id"`
	FishTypeName      string    `json:"fish_type_name"`
	Quality           string    `json:"quality"`
	RemainingWeightKG float64   `json:"remaining_weight_kg"`
	EnteredAt         time.Time `json:"entered_at"`
	ColdStorageName   string    `json:"cold_storage_name"`
	LocationLabel     *string   `json:"location_label"`
	FIFORank          int       `json:"fifo_rank"`
}

type CreateInput struct {
	FishTypeID      string
	ColdStorageID   string
	Quality         string
	InitialWeightKG float64
	EnteredAt       time.Time
	Notes           *string
}

type UpdateQualityInput struct {
	Quality string
	Notes   *string
}

type UpdateLocationInput struct {
	ColdStorageID string
	Notes         *string
}

type MovementInput struct {
	StockBatchID          string
	MovementType          string
	WeightKG              *float64
	PreviousQuality       *string
	NewQuality            *string
	PreviousColdStorageID *string
	NewColdStorageID      *string
	Description           string
}

type ListFilter struct {
	FishTypeID    *string
	Quality       *string
	ColdStorageID *string
	Status        *string
	Sort          string
}

type FIFOFilter struct {
	FishTypeID *string
	Limit      int
	Offset     int
}
