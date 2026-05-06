package stock

import (
	"context"
	"fmt"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Batch, error) {
	var created Batch
	err := s.store.WithTx(ctx, func(repo TxRepository) error {
		batch, err := repo.CreateBatch(ctx, input)
		if err != nil {
			return err
		}

		weight := batch.InitialWeightKG
		if err := repo.CreateMovement(ctx, MovementInput{
			StockBatchID: batch.ID,
			MovementType: MovementIn,
			WeightKG:     &weight,
			Description:  fmt.Sprintf("Stok masuk %.2f kg", weight),
		}); err != nil {
			return err
		}

		created = batch
		return nil
	})
	return created, err
}

func (s *Service) UpdateQuality(ctx context.Context, id string, input UpdateQualityInput) (map[string]string, error) {
	result := map[string]string{}
	err := s.store.WithTx(ctx, func(repo TxRepository) error {
		batch, err := repo.GetBatch(ctx, id)
		if err != nil {
			return err
		}
		if err := repo.UpdateBatchQuality(ctx, id, input.Quality); err != nil {
			return err
		}

		previous := batch.Quality
		if err := repo.CreateMovement(ctx, MovementInput{
			StockBatchID:    id,
			MovementType:    MovementQualityUpdate,
			PreviousQuality: &previous,
			NewQuality:      &input.Quality,
			Description:     notesOrDefault(input.Notes, fmt.Sprintf("Kualitas stok diperbarui dari %s ke %s", previous, input.Quality)),
		}); err != nil {
			return err
		}

		result = map[string]string{
			"id":               id,
			"previous_quality": previous,
			"new_quality":      input.Quality,
		}
		return nil
	})
	return result, err
}

func (s *Service) UpdateLocation(ctx context.Context, id string, input UpdateLocationInput) (map[string]string, error) {
	result := map[string]string{}
	err := s.store.WithTx(ctx, func(repo TxRepository) error {
		batch, err := repo.GetBatch(ctx, id)
		if err != nil {
			return err
		}
		if err := repo.UpdateBatchLocation(ctx, id, input.ColdStorageID); err != nil {
			return err
		}

		previous := batch.ColdStorageID
		if err := repo.CreateMovement(ctx, MovementInput{
			StockBatchID:          id,
			MovementType:          MovementLocationUpdate,
			PreviousColdStorageID: &previous,
			NewColdStorageID:      &input.ColdStorageID,
			Description:           notesOrDefault(input.Notes, "Lokasi stok diperbarui"),
		}); err != nil {
			return err
		}

		result = map[string]string{
			"id":                       id,
			"previous_cold_storage_id": previous,
			"new_cold_storage_id":      input.ColdStorageID,
		}
		return nil
	})
	return result, err
}

func notesOrDefault(notes *string, fallback string) string {
	if notes != nil && *notes != "" {
		return *notes
	}
	return fallback
}
