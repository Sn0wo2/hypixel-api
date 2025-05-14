package hypixel

import (
	"net/http"
	"strings"
)

type PreRequestHook func(request Request) (Response, error)
type CallBack func(request Request, response Response, err error) (Response, error)

type Client struct {
	baseURL        string
	apiKey         string
	httpClient     *http.Client
	rate           *RateLimit
	preRequestHook PreRequestHook
	callBack       CallBack
}

// NewClient creates a new hypixel client
// key is your hypixel api key
//
// https://api.hypixel.net/
func NewClient(key string, rate *RateLimit) *Client {
	return &Client{
		baseURL:    "https://api.hypixel.net/v2/",
		apiKey:     key,
		httpClient: http.DefaultClient,
		rate:       rate,
	}
}

func (c *Client) GetBaseURL() string {
	return c.baseURL
}

func (c *Client) GetAPIKey() string {
	return c.apiKey
}

func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

func (c *Client) GetRate() *RateLimit {
	return c.rate
}

func (c *Client) GetPreRequestHook() PreRequestHook {
	return c.preRequestHook
}

func (c *Client) GetCallBack() CallBack {
	return c.callBack
}

func (c *Client) GetFullPath(path string) string {
	return strings.TrimRight(c.baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *Client) SetAPIKey(key string) {
	c.apiKey = key
}

func (c *Client) SetRate(rate *RateLimit) {
	c.rate = rate
}

func (c *Client) SetPreRequestHook(beforeSend PreRequestHook) {
	c.preRequestHook = beforeSend
}

func (c *Client) SetCallBack(callBack CallBack) {
	c.callBack = callBack
}
