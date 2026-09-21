package calendars

import "testing"

func TestRegistry(t *testing.T) {
	registry := NewRegistry(NewTonalpohualliCASO())
	if _, err := registry.Get(TonalpohualliCASOID); err != nil { t.Fatalf("expected system: %v", err) }
	if _, err := registry.Get("unknown"); err != ErrUnsupportedSystem { t.Fatalf("expected ErrUnsupportedSystem, got %v", err) }
}
