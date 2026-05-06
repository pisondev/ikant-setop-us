package stock

import "testing"

func TestValidateCreateRejectsInvalidRequiredFields(t *testing.T) {
	_, errs := validateCreate(createBody{
		FishTypeID:      "not-uuid",
		ColdStorageID:   "",
		Quality:         "segar",
		InitialWeightKG: 0,
		EnteredAt:       "2026-05-01",
	})
	if len(errs) != 5 {
		t.Fatalf("expected 5 validation errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateCreateAcceptsContractPayload(t *testing.T) {
	notes := " Tangkapan pagi "
	input, errs := validateCreate(createBody{
		FishTypeID:      "7f9d8c5a-4a3b-4f2e-8a21-111111111111",
		ColdStorageID:   "16dfdc88-831f-4e28-b2c7-222222222222",
		Quality:         QualityGood,
		InitialWeightKG: 50,
		EnteredAt:       "2026-05-01T08:00:00Z",
		Notes:           &notes,
	})
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %v", errs)
	}
	if input.Notes == nil || *input.Notes != "Tangkapan pagi" {
		t.Fatalf("expected trimmed notes, got %v", input.Notes)
	}
}

func TestParseRequiredTimeRequiresRFC3339(t *testing.T) {
	_, err := parseRequiredTime("entered_at", "2026-05-01 08:00:00")
	if err == nil {
		t.Fatal("expected time format validation error")
	}
}
