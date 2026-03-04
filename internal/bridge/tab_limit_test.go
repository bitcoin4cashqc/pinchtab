package bridge

import (
	"testing"
	"time"
)

func TestTabLimitError(t *testing.T) {
	err := TabLimitError{Current: 20, Max: 20}
	
	if err.Error() != "tab limit reached (20/20)" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	
	if err.StatusCode() != 429 {
		t.Errorf("expected status code 429, got %d", err.StatusCode())
	}
}

func TestIsTabLimitError(t *testing.T) {
	err := TabLimitError{Current: 20, Max: 20}
	
	if !IsTabLimitError(err) {
		t.Error("IsTabLimitError should return true for TabLimitError")
	}
	
	// Test with wrapped error
	wrappedErr := error(err)
	if !IsTabLimitError(wrappedErr) {
		t.Error("IsTabLimitError should return true for wrapped TabLimitError")
	}
	
	// Test with non-TabLimitError
	otherErr := error(nil)
	if IsTabLimitError(otherErr) {
		t.Error("IsTabLimitError should return false for nil")
	}
}

func TestTabEntry_Metadata(t *testing.T) {
	now := time.Now()
	entry := &TabEntry{
		CreatedAt: now,
		LastUsed:  now,
	}
	
	if entry.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	
	if entry.LastUsed.IsZero() {
		t.Error("LastUsed should not be zero")
	}
	
	// Simulate access
	later := now.Add(5 * time.Second)
	entry.LastUsed = later
	
	if !entry.LastUsed.After(entry.CreatedAt) {
		t.Error("LastUsed should be after CreatedAt after access")
	}
}
