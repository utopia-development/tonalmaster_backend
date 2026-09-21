package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	NewHealthHandler(nil).Health(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK { t.Fatalf("status = %d, want 200", rec.Code) }
}

func TestHealthJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	NewHealthHandler(nil).Health(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if got := rec.Body.String(); got != "{"status":"ok"}\n" { t.Fatalf("body = %q", got) }
}
