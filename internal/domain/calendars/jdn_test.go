package calendars

import (
	"testing"
	"time"
)

func TestGregorianToJDN(t *testing.T) {
	tests := []struct {
		name string
		date string
		want int64
	}{
		{"2026-05-18", "2026-05-18", 2461179},
		{"2000-01-01", "2000-01-01", 2451545},
		{"1970-01-01", "1970-01-01", 2440588},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, err := time.Parse("2006-01-02", tt.date)
			if err != nil {
				t.Fatal(err)
			}
			if got := GregorianToJDN(date); got != tt.want {
				t.Fatalf("GregorianToJDN(%s) = %d, want %d", tt.date, got, tt.want)
			}
		})
	}
}
