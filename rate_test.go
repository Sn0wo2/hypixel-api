package hypixel

import (
	"net/http"
	"testing"
	"time"
)

func TestNewRateLimit(t *testing.T) {
	rl := NewRateLimit()
	if rl.remaining != -1 || !rl.resetAt.IsZero() || rl.trusted {
		t.Errorf("NewRateLimit incorrect initial state: %+v", rl)
	}
}

func TestWaitIfNeeded(t *testing.T) {
	t.Run("no wait when remaining > 0", func(t *testing.T) {
		rl := &RateLimit{remaining: 1}
		start := time.Now()
		rl.WaitIfNeeded()
		if time.Since(start) > 10*time.Millisecond {
			t.Error("Shouldn't wait when remaining > 0")
		}
	})

	t.Run("no wait when reset passed", func(t *testing.T) {
		rl := &RateLimit{
			remaining: 0,
			resetAt:   time.Now().Add(-1 * time.Second),
		}
		start := time.Now()
		rl.WaitIfNeeded()
		if time.Since(start) > 10*time.Millisecond {
			t.Error("Shouldn't wait when reset time has passed")
		}
	})

	t.Run("wait when rate limited", func(t *testing.T) {
		rl := &RateLimit{
			remaining: 0,
			resetAt:   time.Now().Add(100 * time.Millisecond),
		}
		start := time.Now()
		rl.WaitIfNeeded()
		elapsed := time.Since(start)

		if elapsed < 100*time.Millisecond {
			t.Errorf("Should wait at least 100ms, waited %v", elapsed)
		}
	})
}

func TestUpdateFromHeaders(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		headers     map[string]string
		initial     *RateLimit
		wantRem     int
		wantReset   time.Time
		wantTrusted bool
	}{
		{
			name: "initial trusted setup",
			headers: map[string]string{
				"RateLimit-Remaining": "5",
				"RateLimit-Reset":     "300",
			},
			initial:     NewRateLimit(),
			wantRem:     5,
			wantReset:   now.Add(300 * time.Second),
			wantTrusted: true,
		},
		{
			name: "subsequent decrement",
			headers: map[string]string{
				"RateLimit-Remaining": "5",
				"RateLimit-Reset":     "300",
			},
			initial: &RateLimit{
				remaining: 2,
				resetAt:   now.Add(5 * time.Minute),
				trusted:   true,
			},
			wantRem:     1, // should decrement from previous 2
			wantReset:   now.Add(5 * time.Minute),
			wantTrusted: true,
		},
		{
			name: "reset period passed",
			headers: map[string]string{
				"RateLimit-Remaining": "5",
				"RateLimit-Reset":     "300",
			},
			initial: &RateLimit{
				remaining: 0,
				resetAt:   now.Add(-1 * time.Second),
				trusted:   true,
			},
			wantRem:     5,
			wantReset:   now.Add(300 * time.Second),
			wantTrusted: true,
		},
		{
			name:    "ignore invalid headers",
			headers: map[string]string{},
			initial: &RateLimit{
				remaining: 3,
				resetAt:   now.Add(time.Hour),
				trusted:   true,
			},
			wantRem:     3,
			wantReset:   now.Add(time.Hour),
			wantTrusted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := tt.initial
			h := make(http.Header)
			for k, v := range tt.headers {
				h.Set(k, v)
			}

			err := rl.UpdateFromHeaders(h)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if rl.remaining != tt.wantRem || !rl.resetAt.Equal(tt.wantReset) || rl.trusted != tt.wantTrusted {
				t.Errorf("UpdateFromHeaders() = %+v, want remaining %d reset %v trusted %t",
					rl, tt.wantRem, tt.wantReset, tt.wantTrusted)
			}
		})
	}
}

func TestReset(t *testing.T) {
	rl := &RateLimit{
		remaining: 5,
		resetAt:   time.Now().Add(time.Hour),
		trusted:   true,
	}
	rl.Reset()

	if rl.remaining != -1 || !rl.resetAt.IsZero() || rl.trusted {
		t.Errorf("Reset() failed, got: %+v", rl)
	}
}

func TestString(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	rl := &RateLimit{
		remaining: 5,
		resetAt:   fixedTime,
		trusted:   true,
	}

	want := "5 remaining until 2023-01-01T00:00:00Z (trusted:true)"
	if got := rl.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
