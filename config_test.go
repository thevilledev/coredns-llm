package llm

import "testing"

func TestParsePositiveInt(t *testing.T) {
	if _, err := parsePositiveInt("0"); err == nil {
		t.Fatalf("expected error for 0")
	}
	if _, err := parsePositiveInt("-1"); err == nil {
		t.Fatalf("expected error for -1")
	}
	if v, err := parsePositiveInt("10"); err != nil || v != 10 {
		t.Fatalf("expected 10, got %v, err=%v", v, err)
	}
}
