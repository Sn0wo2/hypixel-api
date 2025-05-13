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
