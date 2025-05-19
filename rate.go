package hypixel

import (
	"math"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type RateLimit struct {
	remaining atomic.Int32 // -1 == unknown, >0 == calls left
	resetAt   atomic.Value // holds time.Time
}

func NewRateLimit() *RateLimit {
	r := &RateLimit{}
	r.remaining.Store(-1)
	r.resetAt.Store(time.Time{})
	return r
}

// WaitIfNeeded blocks until rate-limit reset if remaining ≤ 0 and resetAt is in the future.
func (r *RateLimit) WaitIfNeeded() {
	if r.remaining.Load() >= 0 {
		return
	}

	reset := r.resetAt.Load().(time.Time)
	if reset.IsZero() || time.Now().After(reset) {
		return
	}

	// Sleep up to the reset time (capped to hypixel api max reset: 5min)
	sleep := time.Until(reset)
	const maxSleep = 5 * time.Minute
	if sleep > maxSleep {
		sleep = maxSleep
	}
	time.Sleep(sleep)
}

// UpdateFromResponse updates rate limit state based on the HTTP response.
func (r *RateLimit) UpdateFromResponse(resp *http.Response) error {
	resetStr := resp.Header.Get("RateLimit-Reset")
	// always trust reset header if not empty
	if resetStr != "" {
		if secs, err := strconv.Atoi(resetStr); err == nil {
			resetTime := time.Now().Add(time.Duration(secs) * time.Second)
			r.resetAt.Store(resetTime)
		} else {
			return err
		}
	}

	if resp.StatusCode == 429 {
		r.remaining.Store(-1)
		return nil
	}

	remStr := resp.Header.Get("RateLimit-Remaining")
	rem, err := strconv.Atoi(remStr)
	if err != nil {
		return err
	}

	// trust remaining header if status code not 429
	//
	// Example(api return 0 remaining header):
	// {"success":false,"cause":"You have already looked up this player too recently, please try again shortly"}
	// {"success":false,"cause":"Too many requests in the last second","throttle":true}
	if remStr != "" {
		switch {
		case rem > math.MinInt32 && rem <= math.MaxInt32:
			r.remaining.Store(int32(rem))
		default:
			r.remaining.Store(-1)
		}
		return nil
	}

	if r.remaining.Load() > 0 {
		r.remaining.Add(-1)
	}
	return nil
}

// Reset clears all rate-limit state
func (r *RateLimit) Reset() {
	r.remaining.Store(-1)
	r.resetAt.Store(time.Time{})
}

// String impl fmt.Stringer
func (r *RateLimit) String() string {
	reset := r.resetAt.Load().(time.Time)
	return strconv.Itoa(int(r.remaining.Load())) +
		" remaining until " + reset.Format(time.RFC3339)
}
