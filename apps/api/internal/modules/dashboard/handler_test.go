package dashboard

import "testing"

func TestParseLimitDefaultsToTen(t *testing.T) {
	limit, err := parseLimit("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if limit != 10 {
		t.Fatalf("expected default limit 10, got %d", limit)
	}
}

func TestParseLimitRejectsInvalidValue(t *testing.T) {
	_, err := parseLimit("0")
	if err == nil || err.Field != "limit" {
		t.Fatalf("expected limit validation error, got %v", err)
	}
}
