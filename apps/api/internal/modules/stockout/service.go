package stockout

import (
	"context"
	"errors"
	"fmt"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

type InsufficientStockError struct {
	Requested float64
	Available float64
}

func (e InsufficientStockError) Error() string {
	return fmt.Sprintf("requested %.2f kg, but only %.2f kg is available", e.Requested, e.Available)
}

func (s *Service) Create(ctx context.Context, input CreateInput) (StockOut, error) {
	var created StockOut
	err := s.store.WithTx(ctx, func(repo TxRepository) error {
		batches, err := repo.ListAvailableBatchesForUpdate(ctx, input.FishTypeID)
		if err != nil {
			return err
		}

		available := sumRemaining(batches)
		if available+0.000001 < input.TotalWeightKG {
			return InsufficientStockError{Requested: input.TotalWeightKG, Available: available}
		}

		stockOut, err := repo.CreateRecord(ctx, input)
		if err != nil {
			return err
		}

		remainingToTake := input.TotalWeightKG
		items := make([]Item, 0, len(batches))
		for _, batch := range batches {
			if remainingToTake <= 0 {
				break
			}

			taken := minFloat(batch.RemainingWeightKG, remainingToTake)
			newRemaining := roundWeight(batch.RemainingWeightKG - taken)
			status := StatusAvailable
			if newRemaining <= 0 {
				newRemaining = 0
				status = StatusDepleted
			}

			item := Item{
				StockBatchID: batch.ID,
				WeightKG:     roundWeight(taken),
			}
			if err := repo.CreateItem(ctx, stockOut.ID, item); err != nil {
				return err
			}
			if err := repo.UpdateBatchRemaining(ctx, batch.ID, newRemaining, status); err != nil {
				return err
			}

			weight := item.WeightKG
			if err := repo.CreateMovement(ctx, MovementInput{
				StockBatchID: batch.ID,
				MovementType: MovementOut,
				WeightKG:     &weight,
				Description:  fmt.Sprintf("Stok keluar %.2f kg ke %s", weight, input.Destination),
			}); err != nil {
				return err
			}

			items = append(items, item)
			remainingToTake = roundWeight(remainingToTake - taken)
		}

		if remainingToTake > 0.000001 {
			return errors.New("fifo stock deduction did not fulfill requested weight")
		}

		stockOut.FishTypeID = input.FishTypeID
		stockOut.Items = items
		created = stockOut
		return nil
	})
	return created, err
}

func sumRemaining(batches []AvailableBatch) float64 {
	total := 0.0
	for _, batch := range batches {
		total += batch.RemainingWeightKG
	}
	return roundWeight(total)
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func roundWeight(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}
