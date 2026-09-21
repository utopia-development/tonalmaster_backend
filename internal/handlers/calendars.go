package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/utopia-development/tonalmaster_backend/internal/domain/calendars"
)

type CalendarHandler struct {
	registry *calendars.Registry
}

func NewCalendarHandler(registry *calendars.Registry) *CalendarHandler {
	return &CalendarHandler{registry: registry}
}

type calendarDTO struct {
	ID string `json:"id"`
	Name string `json:"nombre"`
}

type conversionDTO struct {
	GregorianDate string `json:"fecha_gregoriana"`
	System string `json:"sistema"`
	JDN int64 `json:"jdn"`
	Result resultDTO `json:"resultado"`
}

type resultDTO struct {
	Trecena int `json:"trecena"`
	Sign string `json:"signo,omitempty"`
	DayNumber int `json:"numero_dia"`
	NightLord string `json:"senor_de_la_noche,omitempty"`
}

func (h *CalendarHandler) List(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, []calendarDTO{{ID: calendars.TonalpohualliCASOID, Name: "Tonalpohualli"}})
}

func (h *CalendarHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	system, err := h.registry.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "calendar_not_found", "calendar system not found")
		return
	}
	writeJSON(w, http.StatusOK, calendarDTO{ID: system.ID(), Name: "Tonalpohualli"})
}

func (h *CalendarHandler) Convert(w http.ResponseWriter, r *http.Request) {
	dateValue := r.URL.Query().Get("date")
	systemID := r.URL.Query().Get("system")
	if dateValue == "" || systemID == "" {
		writeError(w, http.StatusBadRequest, "invalid_request", "date and system are required")
		return
	}
	date, err := time.Parse("2006-01-02", dateValue)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_date", "date must use YYYY-MM-DD")
		return
	}
	system, err := h.registry.Get(systemID)
	if err != nil {
		writeError(w, http.StatusNotFound, "calendar_not_found", "calendar system not found")
		return
	}
	result, err := system.Convert(date)
	if err != nil {
		if errors.Is(err, calendars.ErrUnsupportedSystem) {
			writeError(w, http.StatusNotFound, "calendar_not_found", "calendar system not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "conversion_failed", "calendar conversion failed")
		return
	}
	writeJSON(w, http.StatusOK, conversionDTO{
		GregorianDate: result.Date.Format("2006-01-02"),
		System: result.System,
		JDN: result.JDN,
		Result: resultDTO{Trecena: result.Trecena, Sign: result.Sign, DayNumber: result.DayNumber, NightLord: result.NightLord},
	})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
