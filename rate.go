package hypixel

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RateLimit struct {
	mu        sync.RWMutex
	remaining int       // current remaining requests
	resetAt   time.Time // next reset time
	trusted   bool      // whether trusted initial headers have been received
}

func NewRateLimit() *RateLimit {
	return &RateLimit{
		remaining: -1, // -1 indicates unknown state
	}
}

// WaitIfNeeded blocks until the next reset window if remaining quota is exhausted
func (r *RateLimit) WaitIfNeeded() {
	r.mu.RLock()
	remaining, reset := r.remaining, r.resetAt
	r.mu.RUnlock()

	if remaining != 0 || time.Now().After(reset) {
		return
	}

	if sleep := time.Until(reset); sleep > 0 {
		const maxSleep = 5 * time.Minute // prevent excessive waiting
		if sleep > maxSleep {
			sleep = maxSleep
		}
		time.Sleep(sleep)
	}
}

// UpdateFromHeaders updates rate limits from HTTP headers, auto-consuming on trusted state
func (r *RateLimit) UpdateFromHeaders(h http.Header) error {
	remStr, resetStr := h.Get("RateLimit-Remaining"), h.Get("RateLimit-Reset")
	if remStr == "" || resetStr == "" {
		return nil // ignore non-rate-limited responses
	}

	rem, err := strconv.Atoi(remStr)
	if err != nil {
		return err
	}

	secs, err := strconv.Atoi(resetStr)
	if err != nil {
		return err
	}
	newReset := time.Now().Add(time.Duration(secs) * time.Second)

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	if !r.trusted || now.After(r.resetAt) {
		if rem > 0 { // Only trust headers with positive remaining count
			r.remaining = rem
			r.resetAt = newReset
			r.trusted = true
		}
	} else if r.remaining > 0 {
		r.remaining-- // Local consumption for subsequent requests
	}

	return nil
}

// Reset clears all rate limiting state
func (r *RateLimit) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.remaining = -1
	r.resetAt = time.Time{}
	r.trusted = false
}

// String returns current rate limit status for debugging
func (r *RateLimit) String() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return strconv.Itoa(r.remaining) + " remaining until " +
		r.resetAt.Format(time.RFC3339) + " (trusted:" +
		strconv.FormatBool(r.trusted) + ")"
}
