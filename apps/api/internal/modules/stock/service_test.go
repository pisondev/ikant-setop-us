package stock

import (
	"context"
	"testing"
	"time"
)

func TestCreateCreatesAvailableBatchAndInMovement(t *testing.T) {
	repo := newFakeStore()
	service := NewService(repo)

	created, err := service.Create(context.Background(), CreateInput{
		FishTypeID:      "fish-type-1",
		ColdStorageID:   "storage-1",
		Quality:         QualityGood,
		InitialWeightKG: 50,
		EnteredAt:       time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.RemainingWeightKG != 50 || created.Status != StatusAvailable {
		t.Fatalf("expected available batch with 50 kg, got status=%s remaining=%.2f", created.Status, created.RemainingWeightKG)
	}
	if len(repo.movements) != 1 || repo.movements[0].MovementType != MovementIn {
		t.Fatalf("expected one in movement, got %#v", repo.movements)
	}
}

func TestUpdateQualityCreatesMovementWithPreviousQuality(t *testing.T) {
	repo := newFakeStore()
	repo.batches["batch-1"] = Batch{
		ID:            "batch-1",
		Quality:       QualityGood,
		ColdStorageID: "storage-1",
		Status:        StatusAvailable,
	}
	service := NewService(repo)

	result, err := service.UpdateQuality(context.Background(), "batch-1", UpdateQualityInput{Quality: QualityMedium})
	if err != nil {
		t.Fatalf("UpdateQuality returned error: %v", err)
	}
	if result["previous_quality"] != QualityGood || result["new_quality"] != QualityMedium {
		t.Fatalf("unexpected result: %#v", result)
	}
	if repo.batches["batch-1"].Quality != QualityMedium {
		t.Fatalf("expected batch quality to be updated")
	}
	if len(repo.movements) != 1 || repo.movements[0].MovementType != MovementQualityUpdate {
		t.Fatalf("expected quality movement, got %#v", repo.movements)
	}
}

func TestUpdateLocationCreatesMovementWithPreviousLocation(t *testing.T) {
	repo := newFakeStore()
	repo.batches["batch-1"] = Batch{
		ID:            "batch-1",
		Quality:       QualityGood,
		ColdStorageID: "storage-old",
		Status:        StatusAvailable,
	}
	service := NewService(repo)

	result, err := service.UpdateLocation(context.Background(), "batch-1", UpdateLocationInput{ColdStorageID: "storage-new"})
	if err != nil {
		t.Fatalf("UpdateLocation returned error: %v", err)
	}
	if result["previous_cold_storage_id"] != "storage-old" || result["new_cold_storage_id"] != "storage-new" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if repo.batches["batch-1"].ColdStorageID != "storage-new" {
		t.Fatalf("expected batch location to be updated")
	}
	if len(repo.movements) != 1 || repo.movements[0].MovementType != MovementLocationUpdate {
		t.Fatalf("expected location movement, got %#v", repo.movements)
	}
}

type fakeStore struct {
	batches   map[string]Batch
	movements []MovementInput
	nextID    int
}

func newFakeStore() *fakeStore {
	return &fakeStore{batches: map[string]Batch{}, nextID: 1}
}

func (f *fakeStore) WithTx(ctx context.Context, fn func(TxRepository) error) error {
	snapshot := f.clone()
	if err := fn(f); err != nil {
		*f = *snapshot
		return err
	}
	return nil
}

func (f *fakeStore) CreateBatch(ctx context.Context, input CreateInput) (Batch, error) {
	id := "batch-1"
	if f.nextID > 1 {
		id = "batch-2"
	}
	f.nextID++
	batch := Batch{
		ID:                id,
		FishTypeID:        input.FishTypeID,
		ColdStorageID:     input.ColdStorageID,
		Quality:           input.Quality,
		InitialWeightKG:   input.InitialWeightKG,
		RemainingWeightKG: input.InitialWeightKG,
		EnteredAt:         input.EnteredAt,
		Status:            StatusAvailable,
	}
	f.batches[id] = batch
	return batch, nil
}

func (f *fakeStore) CreateMovement(ctx context.Context, input MovementInput) error {
	f.movements = append(f.movements, input)
	return nil
}

func (f *fakeStore) GetBatch(ctx context.Context, id string) (Batch, error) {
	batch, ok := f.batches[id]
	if !ok {
		return Batch{}, ErrNotFound
	}
	return batch, nil
}

func (f *fakeStore) UpdateBatchLocation(ctx context.Context, id string, coldStorageID string) error {
	batch := f.batches[id]
	batch.ColdStorageID = coldStorageID
	f.batches[id] = batch
	return nil
}

func (f *fakeStore) UpdateBatchQuality(ctx context.Context, id string, quality string) error {
	batch := f.batches[id]
	batch.Quality = quality
	f.batches[id] = batch
	return nil
}

func (f *fakeStore) List(ctx context.Context, filter ListFilter) ([]Detail, error) {
	return nil, nil
}

func (f *fakeStore) GetDetail(ctx context.Context, id string) (Detail, error) {
	return Detail{}, nil
}

func (f *fakeStore) ListFIFO(ctx context.Context, filter FIFOFilter) ([]FIFOItem, error) {
	return nil, nil
}

func (f *fakeStore) clone() *fakeStore {
	clone := &fakeStore{
		batches:   map[string]Batch{},
		movements: append([]MovementInput(nil), f.movements...),
		nextID:    f.nextID,
	}
	for id, batch := range f.batches {
		clone.batches[id] = batch
	}
	return clone
}
