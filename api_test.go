package hypixel

import (
	"net/http"
	"testing"
)

func TestClient_Authentication(t *testing.T) {
	h := http.Header{}
	h.Set("head", "value1")
	c := NewClient("test1", nil)
	if c.AuthHeader(h).Get("API-Key") != "test1" {
		t.Errorf("expected 'test1', got %s", h.Get("API-Key"))
	}
	if h.Get("head") != "value1" {
		t.Errorf("expected 'value1', got %s", h.Get("head"))
	}
}
