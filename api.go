package hypixel

import (
	"bytes"
	"io"
	"net/http"
)

func (c *Client) Send(method string, head http.Header, path string, params *Params, payload ...byte) (*http.Response, error) {
	if method == "" {
		method = http.MethodGet
	}
	full := c.GetFullPath(path)
	if params != nil {
		full = params.String(full)
	}
	req, err := http.NewRequest(method, full,
		func() io.Reader {
			if payload != nil {
				return bytes.NewReader(payload)
			}
			return nil
		}(),
	)
	if err != nil {
		return nil, err
	}
	if head != nil {
		req.Header = head
	}
	c.Rate.WaitIfNeeded()
	rsp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	_ = c.Rate.UpdateRateLimitInfo(rsp.Header)
	return rsp, nil
}

// Authentication Add api key to header
//
// https://api.hypixel.net/#section/Authentication/ApiKey
func (c *Client) Authentication(header ...http.Header) http.Header {
	var h http.Header
	if len(header) == 0 {
		h = http.Header{}
	} else {
		h = header[0]
	}
	h.Set("API-Key", c.APIKey)
	return h
}

// GetPlayerData Data of a specific player, including game stats
//
// https://api.hypixel.net/#tag/Player-Data
// 200 Get player's data
// 400 Some data is missing, this is usually a field.
// 403 Access is forbidden, usually due to an invalid API key being used.
// 429 A request limit has been reached, usually this is due to the limit on the key being reached but can also be triggered by a global throttle.
func (c *Client) GetPlayerData(uuid string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "v2/player", &Params{
		"uuid": uuid,
	})
}
