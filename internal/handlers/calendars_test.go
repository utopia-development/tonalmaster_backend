package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/utopia-development/tonalmaster_backend/internal/domain/calendars"
)

func TestCalendarList(t *testing.T) {
	h := NewCalendarHandler(calendars.NewRegistry(calendars.NewTonalpohualliCASO()))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("status = %d", rec.Code) }
}

func TestCalendarConvert(t *testing.T) {
	h := NewCalendarHandler(calendars.NewRegistry(calendars.NewTonalpohualliCASO()))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars/convert?date=2026-05-18&system=tonalpohualli_caso", nil)
	rec := httptest.NewRecorder()
	h.Convert(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("status = %d", rec.Code) }
}

func TestCalendarConvertInvalidDate(t *testing.T) {
	h := NewCalendarHandler(calendars.NewRegistry(calendars.NewTonalpohualliCASO()))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars/convert?date=bad&system=tonalpohualli_caso", nil)
	rec := httptest.NewRecorder()
	h.Convert(rec, req)
	if rec.Code != http.StatusBadRequest { t.Fatalf("status = %d", rec.Code) }
}

func TestCalendarUnknownSystem(t *testing.T) {
	h := NewCalendarHandler(calendars.NewRegistry(calendars.NewTonalpohualliCASO()))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars/convert?date=2026-05-18&system=unknown", nil)
	rec := httptest.NewRecorder()
	h.Convert(rec, req)
	if rec.Code != http.StatusNotFound { t.Fatalf("status = %d", rec.Code) }
}
