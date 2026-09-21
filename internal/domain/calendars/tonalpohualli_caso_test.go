package calendars

import (
	"testing"
	"time"
)

func TestTonalpohualliCASOJDNAndCycle(t *testing.T) {
	system := NewTonalpohualliCASO(0)
	date := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)

	got, err := system.Convert(date)
	if err != nil {
		t.Fatal(err)
	}
	if got.JDN != 2461179 {
		t.Fatalf("JDN = %d, want 2461179", got.JDN)
	}
	if got.Trecena < 1 || got.Trecena > 13 {
		t.Fatalf("Trecena out of range: %d", got.Trecena)
	}
	if got.DayNumber < 1 || got.DayNumber > 20 {
		t.Fatalf("DayNumber out of range: %d", got.DayNumber)
	}
}
