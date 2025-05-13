package hypixel

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	rateLimit := &RateLimit{}

	client := NewClient(apiKey, rateLimit)

	if client.APIKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.APIKey)
	}

	if client.BaseURL != "https://api.hypixel.net/v2/" {
		t.Errorf("Expected base URL %s, got %s", "https://api.hypixel.net/v2/", client.BaseURL)
	}

	if client.HTTP != http.DefaultClient {
		t.Error("Expected default HTTP client")
	}

	if client.Rate != rateLimit {
		t.Error("Rate limit not set correctly")
	}
}

func TestGetters(t *testing.T) {
	customHTTP := &http.Client{}
	client := &Client{
		BaseURL: "https://custom.url/",
		APIKey:  "custom-key",
		HTTP:    customHTTP,
		Rate:    &RateLimit{},
	}

	t.Run("GetBaseURL", func(t *testing.T) {
		if got := client.GetBaseURL(); got != "https://custom.url/" {
			t.Errorf("GetBaseURL() = %v, want %v", got, "https://custom.url/")
		}
	})

	t.Run("GetAPIKey", func(t *testing.T) {
		if got := client.GetAPIKey(); got != "custom-key" {
			t.Errorf("GetAPIKey() = %v, want %v", got, "custom-key")
		}
	})

	t.Run("GetHTTPClient", func(t *testing.T) {
		if got := client.GetHTTPClient(); got != customHTTP {
			t.Errorf("GetHTTPClient() = %v, want %v", got, customHTTP)
		}
	})
}

func TestSetters(t *testing.T) {
	client := NewClient("initial-key", &RateLimit{})

	t.Run("SetBaseURL", func(t *testing.T) {
		newURL := "https://new.url/"
		client.SetBaseURL(newURL)
		if client.BaseURL != newURL {
			t.Errorf("SetBaseURL() failed, got %v, want %v", client.BaseURL, newURL)
		}
	})

	t.Run("SetAPIKey", func(t *testing.T) {
		newKey := "new-key"
		client.SetAPIKey(newKey)
		if client.APIKey != newKey {
			t.Errorf("SetAPIKey() failed, got %v, want %v", client.APIKey, newKey)
		}
	})

	t.Run("SetHTTPClient", func(t *testing.T) {
		newClient := &http.Client{}
		client.SetHTTPClient(newClient)
		if client.HTTP != newClient {
			t.Errorf("SetHTTPClient() failed, got %v, want %v", client.HTTP, newClient)
		}
	})
}

func TestGetFullPath(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		path     string
		expected string
	}{
		{
			name:     "baseURL with trailing slash",
			baseURL:  "https://api.hypixel.net/v2/",
			path:     "skyblock/auctions",
			expected: "https://api.hypixel.net/v2/skyblock/auctions",
		},
		{
			name:     "baseURL without trailing slash",
			baseURL:  "https://api.hypixel.net/v2",
			path:     "skyblock/auctions",
			expected: "https://api.hypixel.net/v2/skyblock/auctions",
		},
		{
			name:     "empty path",
			baseURL:  "https://api.hypixel.net/v2/",
			path:     "",
			expected: "https://api.hypixel.net/v2/",
		},
		{
			name:     "path with leading slash",
			baseURL:  "https://api.hypixel.net/v2/",
			path:     "/skyblock/auctions",
			expected: "https://api.hypixel.net/v2/skyblock/auctions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{BaseURL: tt.baseURL}
			result := client.GetFullPath(tt.path)
			if result != tt.expected {
				t.Errorf("GetFullPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}
