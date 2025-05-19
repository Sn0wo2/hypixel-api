package hypixel

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type RateLimit struct {
	remaining atomic.Int32  // -1 == unknown, >0 == calls left
	resetAt   atomic.Value  // holds time.Time
	mu        sync.Mutex    // protects waitCh
	waitCh    chan struct{} // closed when reset time is reached
}

func NewRateLimit() *RateLimit {
	r := &RateLimit{}
	r.remaining.Store(-1)
	r.resetAt.Store(time.Time{})
	return r
}

// WaitIfNeeded blocks until rate-limit reset if remaining ≤ 0 and resetAt is in the future.
func (r *RateLimit) WaitIfNeeded() {
	for {
		r.mu.Lock()
		rem := r.remaining.Load()
		reset := r.resetAt.Load().(time.Time)

		if rem > 0 || reset.IsZero() || time.Now().After(reset) {
			r.mu.Unlock()
			return
		}

		if r.waitCh == nil {
			ch := make(chan struct{})
			r.waitCh = ch
			go func(ch chan struct{}, reset time.Time) {
				sleep := time.Until(reset)
				if m := 5 * time.Minute; sleep > m { // m: max hypixel api reset cd
					sleep = m
				}
				time.Sleep(sleep)
				r.mu.Lock()
				close(ch)
				r.waitCh = nil
				r.mu.Unlock()
			}(ch, reset)
		}

		ch := r.waitCh
		r.mu.Unlock()

		<-ch
	}
}

// UpdateFromResponse updates rate limit state based on the HTTP response.
func (r *RateLimit) UpdateFromResponse(resp *http.Response) error {
	if resetStr := resp.Header.Get("RateLimit-Reset"); resetStr != "" {
		if secs, err := strconv.Atoi(resetStr); err == nil {
			r.resetAt.Store(time.Now().Add(time.Duration(secs) * time.Second))
		} else {
			return err
		}
	}

	remStr := resp.Header.Get("RateLimit-Remaining")
	if remStr == "" {
		return nil
	}
	// when status code 429, remaining is fake(0)
	if resp.StatusCode == 429 && remStr == "0" {
		r.remaining.Store(-1)
		return nil
	}
	rem, err := strconv.Atoi(remStr)
	if err != nil {
		return err
	}
	// code ql
	if rem > math.MinInt32 && rem <= math.MaxInt32 {
		r.remaining.Store(int32(rem))
	} else {
		// fallback
		r.remaining.Store(-1)
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
