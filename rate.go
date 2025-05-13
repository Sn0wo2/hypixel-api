package hypixel

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimit tracks API call quota and reset time.
type RateLimit struct {
	mu        sync.RWMutex
	remaining int
	resetAt   time.Time
}

// NewRateLimit returns an uninitialized RateLimit (no blocking until headers seen).
func NewRateLimit() *RateLimit {
	return &RateLimit{remaining: -1}
}

// WaitIfNeeded blocks if quota is exhausted until reset time (max 5m).
func (r *RateLimit) WaitIfNeeded() {
	now := time.Now()

	r.mu.RLock()
	rem, reset := r.remaining, r.resetAt
	r.mu.RUnlock()

	if rem < 0 || rem > 0 || !now.Before(reset) {
		return
	}
	wait := reset.Sub(now)
	if wait > 5*time.Minute {
		wait = 5 * time.Minute
	}
	time.Sleep(wait)
}

// UpdateRateLimitInfo parses headers to update quota and reset time.
func (r *RateLimit) UpdateRateLimitInfo(h http.Header) error {
	remStr, resetStr := h.Get("RateLimit-Remaining"), h.Get("RateLimit-Reset")
	if remStr == "" || resetStr == "" {
		return nil
	}
	rem, err := strconv.Atoi(remStr)
	if err != nil {
		return err
	}
	secs, err := strconv.Atoi(resetStr)
	if err != nil {
		return err
	}

	r.mu.Lock()
	r.remaining = rem
	r.resetAt = time.Now().Add(time.Duration(secs) * time.Second)
	r.mu.Unlock()
	return nil
}

// Remaining returns the current remaining request count.
func (r *RateLimit) Remaining() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.remaining
}

// ResetAt returns the time when quota resets.
func (r *RateLimit) ResetAt() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.resetAt
}

func (r *RateLimit) Reset() {
	r.mu.Lock()
	r.remaining = -1
	r.mu.Unlock()
}
