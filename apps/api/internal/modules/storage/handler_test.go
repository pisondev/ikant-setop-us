package storage

import "testing"

func TestValidateCreateRequiresName(t *testing.T) {
	_, errs := validateCreate(createBody{Name: ""})
	if len(errs) != 1 || errs[0].Field != "name" {
		t.Fatalf("expected name validation error, got %v", errs)
	}
}

func TestValidateCreateCleansOptionalFields(t *testing.T) {
	location := "  Zona A - Rak 1  "
	emptyDescription := " "

	input, errs := validateCreate(createBody{
		Name:          " Cold Storage A ",
		LocationLabel: &location,
		Description:   &emptyDescription,
	})
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %v", errs)
	}
	if input.Name != "Cold Storage A" {
		t.Fatalf("expected trimmed name, got %q", input.Name)
	}
	if input.LocationLabel == nil || *input.LocationLabel != "Zona A - Rak 1" {
		t.Fatalf("expected trimmed location label, got %v", input.LocationLabel)
	}
	if input.Description != nil {
		t.Fatalf("expected empty description to become nil, got %v", input.Description)
	}
}
