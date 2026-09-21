package calendars

import (
	"testing"
	"time"
)

func TestMod(t *testing.T) {
	tests := []struct{ value, divisor int64; want int }{
		{5, 3, 2},
		{-1, 3, 2},
		{-260, 260, 0},
	}
	for _, tt := range tests {
		if got := mod(tt.value, tt.divisor); got != tt.want {
			t.Fatalf("mod(%d,%d)=%d, want %d", tt.value, tt.divisor, got, tt.want)
		}
	}
}

func TestTonalpohualliCASO20260518(t *testing.T) {
	system := NewTonalpohualliCASO()
	got, err := system.Convert(time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC))
	if err != nil { t.Fatal(err) }
	if got.JDN != 2461179 { t.Fatalf("JDN = %d, want 2461179", got.JDN) }
	if got.Trecena != 1 { t.Fatalf("Trecena = %d, want 1", got.Trecena) }
	if got.DayNumber != 12 { t.Fatalf("DayNumber = %d, want 12", got.DayNumber) }
	if got.Sign != "Cozcacuauhtli" { t.Fatalf("Sign = %q, want Cozcacuauhtli", got.Sign) }
}

func TestTonalpohualliCASOCycle(t *testing.T) {
	system := NewTonalpohualliCASO()
	date := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)
	got, _ := system.Convert(date)
	next, _ := system.Convert(date.AddDate(0, 0, 260))
	if got.Sign != next.Sign || got.DayNumber != next.DayNumber || got.Trecena != next.Trecena {
		t.Fatalf("260-day cycle did not repeat: got=%+v next=%+v", got, next)
	}
}
