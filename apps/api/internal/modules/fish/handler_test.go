package fish

import "testing"

func TestValidateCreateTrimsNameAndOptionalFields(t *testing.T) {
	image := "  /images/fish/tuna.png  "
	emptyDescription := "   "

	input, errs := validateCreate(createBody{
		Name:        "  Tuna  ",
		ImageURL:    &image,
		Description: &emptyDescription,
	})
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %v", errs)
	}
	if input.Name != "Tuna" {
		t.Fatalf("expected trimmed name, got %q", input.Name)
	}
	if input.ImageURL == nil || *input.ImageURL != "/images/fish/tuna.png" {
		t.Fatalf("expected trimmed image url, got %v", input.ImageURL)
	}
	if input.Description != nil {
		t.Fatalf("expected empty description to become nil, got %v", input.Description)
	}
}

func TestValidateCreateRequiresName(t *testing.T) {
	_, errs := validateCreate(createBody{Name: "   "})
	if len(errs) != 1 || errs[0].Field != "name" {
		t.Fatalf("expected name validation error, got %v", errs)
	}
}
