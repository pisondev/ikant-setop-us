package stockout

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCreateConsumesFIFOAcrossBatches(t *testing.T) {
	ctx := context.Background()
	fishTypeID := "fish-type-1"
	repo := newFakeStore([]AvailableBatch{
		{
			ID:                "batch-old",
			FishTypeID:        fishTypeID,
			RemainingWeightKG: 25,
			Status:            StatusAvailable,
			EnteredAt:         time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
		},
		{
			ID:                "batch-new",
			FishTypeID:        fishTypeID,
			RemainingWeightKG: 50,
			Status:            StatusAvailable,
			EnteredAt:         time.Date(2026, 5, 2, 8, 0, 0, 0, time.UTC),
		},
	})
	service := NewService(repo)

	out, err := service.Create(ctx, CreateInput{
		FishTypeID:    fishTypeID,
		TotalWeightKG: 40,
		Destination:   "Restoran Laut Makassar",
		OutAt:         time.Date(2026, 5, 3, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if len(out.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(out.Items))
	}
	assertItem(t, out.Items[0], "batch-old", 25)
	assertItem(t, out.Items[1], "batch-new", 15)
	assertBatch(t, repo.batches["batch-old"], 0, StatusDepleted)
	assertBatch(t, repo.batches["batch-new"], 35, StatusAvailable)
	if len(repo.movements) != 2 {
		t.Fatalf("expected 2 movements, got %d", len(repo.movements))
	}
}

func TestCreateDepletesExactSingleBatch(t *testing.T) {
	ctx := context.Background()
	fishTypeID := "fish-type-1"
	repo := newFakeStore([]AvailableBatch{
		{ID: "batch-1", FishTypeID: fishTypeID, RemainingWeightKG: 30, Status: StatusAvailable},
	})
	service := NewService(repo)

	out, err := service.Create(ctx, CreateInput{
		FishTypeID:    fishTypeID,
		TotalWeightKG: 30,
		Destination:   "Restoran Laut Makassar",
		OutAt:         time.Date(2026, 5, 3, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if len(out.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(out.Items))
	}
	assertBatch(t, repo.batches["batch-1"], 0, StatusDepleted)
}

func TestCreateReturnsInsufficientStockAndRollsBack(t *testing.T) {
	ctx := context.Background()
	fishTypeID := "fish-type-1"
	repo := newFakeStore([]AvailableBatch{
		{ID: "batch-1", FishTypeID: fishTypeID, RemainingWeightKG: 30, Status: StatusAvailable},
	})
	service := NewService(repo)

	_, err := service.Create(ctx, CreateInput{
		FishTypeID:    fishTypeID,
		TotalWeightKG: 40,
		Destination:   "Restoran Laut Makassar",
		OutAt:         time.Date(2026, 5, 3, 12, 0, 0, 0, time.UTC),
	})
	var insufficient InsufficientStockError
	if !errors.As(err, &insufficient) {
		t.Fatalf("expected InsufficientStockError, got %T: %v", err, err)
	}

	assertBatch(t, repo.batches["batch-1"], 30, StatusAvailable)
	if len(repo.stockOuts) != 0 || len(repo.items) != 0 || len(repo.movements) != 0 {
		t.Fatalf("expected transaction rollback, got outs=%d items=%d movements=%d", len(repo.stockOuts), len(repo.items), len(repo.movements))
	}
}

type fakeStore struct {
	batches   map[string]AvailableBatch
	ordered   []string
	stockOuts []StockOut
	items     []Item
	movements []MovementInput
	nextID    int
}

func newFakeStore(batches []AvailableBatch) *fakeStore {
	repo := &fakeStore{batches: map[string]AvailableBatch{}, nextID: 1}
	for _, batch := range batches {
		repo.batches[batch.ID] = batch
		repo.ordered = append(repo.ordered, batch.ID)
	}
	return repo
}

func (f *fakeStore) WithTx(ctx context.Context, fn func(TxRepository) error) error {
	snapshot := f.clone()
	if err := fn(f); err != nil {
		*f = *snapshot
		return err
	}
	return nil
}

func (f *fakeStore) ListAvailableBatchesForUpdate(ctx context.Context, fishTypeID string) ([]AvailableBatch, error) {
	items := []AvailableBatch{}
	for _, id := range f.ordered {
		batch := f.batches[id]
		if batch.FishTypeID == fishTypeID && batch.Status == StatusAvailable && batch.RemainingWeightKG > 0 {
			items = append(items, batch)
		}
	}
	return items, nil
}

func (f *fakeStore) CreateRecord(ctx context.Context, input CreateInput) (StockOut, error) {
	out := StockOut{
		ID:            "out-1",
		FishTypeID:    input.FishTypeID,
		Destination:   input.Destination,
		TotalWeightKG: input.TotalWeightKG,
		OutAt:         input.OutAt,
		Notes:         input.Notes,
		CreatedAt:     time.Now().UTC(),
	}
	f.stockOuts = append(f.stockOuts, out)
	return out, nil
}

func (f *fakeStore) CreateItem(ctx context.Context, stockOutID string, item Item) error {
	f.items = append(f.items, item)
	return nil
}

func (f *fakeStore) UpdateBatchRemaining(ctx context.Context, id string, remainingWeightKG float64, status string) error {
	batch := f.batches[id]
	batch.RemainingWeightKG = remainingWeightKG
	batch.Status = status
	f.batches[id] = batch
	return nil
}

func (f *fakeStore) CreateMovement(ctx context.Context, input MovementInput) error {
	f.movements = append(f.movements, input)
	return nil
}

func (f *fakeStore) List(ctx context.Context, filter ListFilter) ([]StockOut, error) {
	return nil, nil
}

func (f *fakeStore) clone() *fakeStore {
	clone := &fakeStore{
		batches:   map[string]AvailableBatch{},
		ordered:   append([]string(nil), f.ordered...),
		stockOuts: append([]StockOut(nil), f.stockOuts...),
		items:     append([]Item(nil), f.items...),
		movements: append([]MovementInput(nil), f.movements...),
		nextID:    f.nextID,
	}
	for id, batch := range f.batches {
		clone.batches[id] = batch
	}
	return clone
}

func assertItem(t *testing.T, item Item, batchID string, weight float64) {
	t.Helper()
	if item.StockBatchID != batchID {
		t.Fatalf("expected item stock batch %q, got %q", batchID, item.StockBatchID)
	}
	if item.WeightKG != weight {
		t.Fatalf("expected item weight %.2f, got %.2f", weight, item.WeightKG)
	}
}

func assertBatch(t *testing.T, batch AvailableBatch, remaining float64, status string) {
	t.Helper()
	if batch.RemainingWeightKG != remaining {
		t.Fatalf("expected remaining %.2f for %s, got %.2f", remaining, batch.ID, batch.RemainingWeightKG)
	}
	if batch.Status != status {
		t.Fatalf("expected status %q for %s, got %q", status, batch.ID, batch.Status)
	}
}
