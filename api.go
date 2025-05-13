package hypixel

import (
	"bytes"
	"io"
	"net/http"
)

// Send Hypixel API HTTP Request
// If you want bypass rate limit, use c.Rate.Reset()
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
	if c.Rate != nil {
		c.Rate.WaitIfNeeded()
	}
	rsp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	if c.Rate != nil {
		_ = c.Rate.UpdateFromHeaders(rsp.Header)
	}
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
	return c.Send(http.MethodGet, c.Authentication(), "player", &Params{
		"uuid": uuid,
	})
}

// GetRecentGames The recently played games of a specific player
//
// https://api.hypixel.net/#tag/Player-Data/paths/~1v2~1recentgames/get
// 200 Get player's recent game
// 400 Some data is missing, this is usually a field.
// 403 Access is forbidden, usually due to an invalid API key being used.
// 422 Some data provided is invalid.
// 429 A request limit has been reached, usually this is due to the limit on the key being reached but can also be triggered by a global throttle.
func (c *Client) GetRecentGames(uuid string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "recentgames", &Params{
		"uuid": uuid,
	})
}

// GetStatus The current online status of a specific player
//
// https://api.hypixel.net/#tag/Player-Data/paths/~1v2~1status/get
// 200 Get player status
// 400 Some data is missing, this is usually a field.
// 403 Access is forbidden, usually due to an invalid API key being used.
// 429 A request limit has been reached, usually this is due to the limit on the key being reached but can also be triggered by a global throttle.
func (c *Client) GetStatus(uuid string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "status", &Params{
		"uuid": uuid,
	})
}

// GetGuild Retrieve a Guild by a player, id, or name
//
// https://api.hypixel.net/#tag/Player-Data/paths/~1v2~1guild/get
// 200 Get guild information
// 400 Some data is missing, this is usually a field.
// 403 Access is forbidden, usually due to an invalid API key being used.
// 429 A request limit has been reached, usually this is due to the limit on the key being reached but can also be triggered by a global throttle.
func (c *Client) GetGuild(id, player, name string) (*http.Response, error) {
	return c.Send(http.MethodGet, c.Authentication(), "guild", &Params{
		"id":     id,
		"player": player,
		"name":   name,
	})
}
