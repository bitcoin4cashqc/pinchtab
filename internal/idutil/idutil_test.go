package idutil

import "testing"

func TestStripTabPrefix(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"tab_ABC123", "ABC123"},
		{"tab_A25658CE1BA82659EBE9C93C46CEE63A", "A25658CE1BA82659EBE9C93C46CEE63A"},
		{"rawid", "rawid"},
		{"tab_", "tab_"}, // just prefix with nothing after — no strip
		{"", ""},
	}
	for _, tt := range tests {
		got := StripTabPrefix(tt.input)
		if got != tt.want {
			t.Errorf("StripTabPrefix(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestTabIDFromCDPTarget_Semantic(t *testing.T) {
	m := NewManager()
	cdpID := "A25658CE1BA82659EBE9C93C46CEE63A"
	got := m.TabIDFromCDPTarget(cdpID)
	want := "tab_" + cdpID
	if got != want {
		t.Errorf("TabIDFromCDPTarget(%q) = %q, want %q", cdpID, got, want)
	}
}

func TestStripTabPrefix_RoundTrip(t *testing.T) {
	m := NewManager()
	cdpID := "TARGET123"
	semantic := m.TabIDFromCDPTarget(cdpID)
	back := StripTabPrefix(semantic)
	if back != cdpID {
		t.Errorf("round-trip failed: %q → %q → %q", cdpID, semantic, back)
	}
}
