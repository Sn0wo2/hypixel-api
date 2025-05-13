package hypixel

import (
	"net/http"
	"strings"
)

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
	Rate    *RateLimit
}

// NewClient creates a new hypixel client
// key is your hypixel api key
//
// https://api.hypixel.net/
func NewClient(key string, rate *RateLimit) *Client {
	return &Client{
		BaseURL: "https://api.hypixel.net/v2/",
		APIKey:  key,
		HTTP:    http.DefaultClient,
		Rate:    rate,
	}
}

func (c *Client) GetBaseURL() string {
	return c.BaseURL
}

func (c *Client) GetAPIKey() string {
	return c.APIKey
}

func (c *Client) GetHTTPClient() *http.Client {
	return c.HTTP
}

func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.HTTP = client
}

func (c *Client) SetAPIKey(key string) {
	c.APIKey = key
}

func (c *Client) GetFullPath(path string) string {
	sb := strings.Builder{}
	sb.Grow(len(c.BaseURL) + len(path) + 1)
	sb.WriteString(c.BaseURL)
	if !strings.HasSuffix(c.BaseURL, "/") {
		sb.WriteString("/")
	}
	sb.WriteString(path)
	return sb.String()
}
