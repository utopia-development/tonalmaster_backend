package calendars

import (
	"testing"
	"time"
)

func TestTonalpohualliCASO20260518(t *testing.T) {
	system := NewTonalpohualliCASO()
	date := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)

	got, err := system.Convert(date)
	if err != nil {
		t.Fatal(err)
	}

	if got.JDN != 2461179 {
		t.Fatalf("JDN = %d, want 2461179", got.JDN)
	}
	if got.Trecena != 1 {
		t.Fatalf("Trecena = %d, want 1", got.Trecena)
	}
	if got.DayNumber != 12 {
		t.Fatalf("DayNumber = %d, want 12", got.DayNumber)
	}
	if got.Sign != "Cozcacuauhtli" {
		t.Fatalf("Sign = %q, want Cozcacuauhtli", got.Sign)
	}
}

func TestTonalpohualliCASOCycle(t *testing.T) {
	system := NewTonalpohualliCASO()
	date := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)

	got, err := system.Convert(date)
	if err != nil {
		t.Fatal(err)
	}
	next, err := system.Convert(date.AddDate(0, 0, 260))
	if err != nil {
		t.Fatal(err)
	}

	if got.Sign != next.Sign || got.DayNumber != next.DayNumber || got.Trecena != next.Trecena {
		t.Fatalf("260-day cycle did not repeat: got=%+v next=%+v", got, next)
	}
}
