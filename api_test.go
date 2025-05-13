package hypixel

import (
	"net/http"
	"testing"
)

func TestClient_Authentication(t *testing.T) {
	h := http.Header{}
	if NewClient("test1", nil).Authentication(h).Get("API-Key") != "test1" {
		t.Errorf("expected 'test1', got %s", h.Get("API-Key"))
	}
}
