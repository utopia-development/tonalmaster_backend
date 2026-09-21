package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/utopia-development/tonalmaster_backend/internal/domain/calendars"
)

func calendarHandlerForTest() *CalendarHandler {
	return NewCalendarHandler(calendars.NewRegistry(calendars.NewTonalpohualliCASO()))
}

func TestCalendarList(t *testing.T) {
	rec := httptest.NewRecorder()
	calendarHandlerForTest().List(rec, httptest.NewRequest(http.MethodGet, "/api/v1/calendars", nil))
	if rec.Code != http.StatusOK { t.Fatalf("status = %d, want 200", rec.Code) }

	var got []calendarDTO
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil { t.Fatal(err) }
	if len(got) != 1 || got[0].ID != calendars.TonalpohualliCASOID { t.Fatalf("unexpected response: %+v", got) }
}

func TestCalendarGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars/"+calendars.TonalpohualliCASOID, nil)
	req.SetPathValue("id", calendars.TonalpohualliCASOID)
	rec := httptest.NewRecorder()
	calendarHandlerForTest().Get(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("status = %d, want 200", rec.Code) }
}

func TestCalendarGetUnknown(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars/unknown", nil)
	req.SetPathValue("id", "unknown")
	rec := httptest.NewRecorder()
	calendarHandlerForTest().Get(rec, req)
	if rec.Code != http.StatusNotFound { t.Fatalf("status = %d, want 404", rec.Code) }
}

func TestCalendarConvert(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars/convert?date=2026-05-18&system=tonalpohualli_caso", nil)
	rec := httptest.NewRecorder()
	calendarHandlerForTest().Convert(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("status = %d, want 200", rec.Code) }

	var got conversionDTO
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil { t.Fatal(err) }
	if got.JDN != 2461179 || got.Result.DayNumber != 12 || got.Result.Sign != "Cozcacuauhtli" {
		t.Fatalf("unexpected conversion: %+v", got)
	}
}

func TestCalendarConvertMissingQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars/convert", nil)
	rec := httptest.NewRecorder()
	calendarHandlerForTest().Convert(rec, req)
	if rec.Code != http.StatusBadRequest { t.Fatalf("status = %d, want 400", rec.Code) }
}

func TestCalendarConvertInvalidDate(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars/convert?date=2026-99-99&system=tonalpohualli_caso", nil)
	rec := httptest.NewRecorder()
	calendarHandlerForTest().Convert(rec, req)
	if rec.Code != http.StatusBadRequest { t.Fatalf("status = %d, want 400", rec.Code) }
}

func TestCalendarConvertUnknownSystem(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calendars/convert?date=2026-05-18&system=unknown", nil)
	rec := httptest.NewRecorder()
	calendarHandlerForTest().Convert(rec, req)
	if rec.Code != http.StatusNotFound { t.Fatalf("status = %d, want 404", rec.Code) }
}
