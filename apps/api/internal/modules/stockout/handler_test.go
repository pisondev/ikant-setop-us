package stockout

import "testing"

func TestValidateCreateRejectsInvalidRequiredFields(t *testing.T) {
	_, errs := validateCreate(createBody{
		FishTypeID:    "not-uuid",
		TotalWeightKG: 0,
		Destination:   " ",
		OutAt:         "2026-05-01",
	})
	if len(errs) != 4 {
		t.Fatalf("expected 4 validation errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateCreateAcceptsContractPayload(t *testing.T) {
	notes := " Pengeluaran untuk pesanan "
	input, errs := validateCreate(createBody{
		FishTypeID:    "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
		TotalWeightKG: 40,
		Destination:   " Restoran Laut Makassar ",
		OutAt:         "2026-05-01T12:00:00Z",
		Notes:         &notes,
	})
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %v", errs)
	}
	if input.Destination != "Restoran Laut Makassar" {
		t.Fatalf("expected trimmed destination, got %q", input.Destination)
	}
	if input.Notes == nil || *input.Notes != "Pengeluaran untuk pesanan" {
		t.Fatalf("expected trimmed notes, got %v", input.Notes)
	}
}
