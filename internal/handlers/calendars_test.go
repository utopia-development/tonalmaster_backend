package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/utopia-development/tonalmaster_backend/internal/domain/calendars"
)

func TestCalendarAPIContractErrors(t *testing.T) {
	h := NewCalendarHandler(calendars.NewRegistry(calendars.NewTonalpohualliCASO()))

	tests := []struct {
		name, path, wantCode, wantMessage string
	}{
		{"missing parameters", "/api/v1/calendars/convert", "invalid_request", "date and system are required"},
		{"invalid date", "/api/v1/calendars/convert?date=2026-99-99&system=tonalpohualli_caso", "invalid_date", "date must use YYYY-MM-DD"},
		{"unknown calendar", "/api/v1/calendars/convert?date=2026-05-18&system=unknown", "calendar_not_found", "calendar system not found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.Convert(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if rec.Code == http.StatusOK { t.Fatalf("expected error, got 200") }
			if got := rec.Header().Get("Content-Type"); got != "application/json" { t.Fatalf("content type = %q", got) }
			body := rec.Body.String()
			if !containsAll(body, tt.wantCode, tt.wantMessage) { t.Fatalf("body = %q", body) }
		})
	}
}

func TestCalendarGetPathValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars/tonalpohualli_caso", nil)
	req.SetPathValue("id", calendars.TonalpohualliCASOID)
	rec := httptest.NewRecorder()
	NewCalendarHandler(calendars.NewRegistry(calendars.NewTonalpohualliCASO())).Get(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("status = %d, want 200", rec.Code) }
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !contains(value, part) { return false }
	}
	return true
}

func contains(value, part string) bool {
	return len(part) == 0 || indexOf(value, part) >= 0
}

func indexOf(value, part string) int {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part { return i }
	}
	return -1
}
