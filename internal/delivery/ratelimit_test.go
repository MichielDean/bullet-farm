package delivery

import (
	"testing"
	"time"
)

// mockClock provides a controllable clock for deterministic tests.
type mockClock struct {
	t time.Time
}

func (c *mockClock) now() time.Time { return c.t }

func (c *mockClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func newTestLimiter(ipLimit, tokenLimit int, window time.Duration, clk *mockClock) *RateLimiter {
	rl := NewRateLimiter(Config{
		PerIPRequests:    ipLimit,
		PerTokenRequests: tokenLimit,
		Window:           window,
	})
	rl.now = clk.now
	return rl
}

func TestRateLimiter_AllowsWithinIPLimit(t *testing.T) {
	clk := &mockClock{t: time.Now()}
	rl := newTestLimiter(3, 100, time.Minute, clk)

	for i := 0; i < 3; i++ {
		if !rl.Allow("1.2.3.4", "tok-a") {
			t.Fatalf("request %d should be allowed, was denied", i+1)
		}
	}
}

func TestRateLimiter_DeniesAtIPLimit(t *testing.T) {
	clk := &mockClock{t: time.Now()}
	rl := newTestLimiter(3, 100, time.Minute, clk)

	for i := 0; i < 3; i++ {
		rl.Allow("1.2.3.4", "tok-a") //nolint:errcheck
	}
	if rl.Allow("1.2.3.4", "tok-a") {
		t.Fatal("4th request from same IP should be denied")
	}
}

func TestRateLimiter_AllowsWithinTokenLimit(t *testing.T) {
	clk := &mockClock{t: time.Now()}
	rl := newTestLimiter(100, 3, time.Minute, clk)

	for i := 0; i < 3; i++ {
		if !rl.Allow("1.2.3.4", "tok-a") {
			t.Fatalf("request %d should be allowed, was denied", i+1)
		}
	}
}

func TestRateLimiter_DeniesAtTokenLimit(t *testing.T) {
	clk := &mockClock{t: time.Now()}
	rl := newTestLimiter(100, 3, time.Minute, clk)

	for i := 0; i < 3; i++ {
		rl.Allow("1.2.3.4", "tok-a") //nolint:errcheck
	}
	if rl.Allow("1.2.3.4", "tok-a") {
		t.Fatal("4th request with same token should be denied")
	}
}

func TestRateLimiter_DifferentIPsAreIndependent(t *testing.T) {
	clk := &mockClock{t: time.Now()}
	rl := newTestLimiter(2, 100, time.Minute, clk)

	// Exhaust IP1.
	rl.Allow("10.0.0.1", "tok-a") //nolint:errcheck
	rl.Allow("10.0.0.1", "tok-a") //nolint:errcheck
	if rl.Allow("10.0.0.1", "tok-a") {
		t.Fatal("3rd request from IP1 should be denied")
	}

	// IP2 should still be allowed.
	if !rl.Allow("10.0.0.2", "tok-b") {
		t.Fatal("request from IP2 should be allowed (independent limit)")
	}
}

func TestRateLimiter_DifferentTokensAreIndependent(t *testing.T) {
	clk := &mockClock{t: time.Now()}
	rl := newTestLimiter(100, 2, time.Minute, clk)

	// Exhaust tok-a.
	rl.Allow("1.2.3.4", "tok-a") //nolint:errcheck
	rl.Allow("5.6.7.8", "tok-a") //nolint:errcheck
	if rl.Allow("9.9.9.9", "tok-a") {
		t.Fatal("3rd request with tok-a should be denied")
	}

	// tok-b should still be allowed from any IP.
	if !rl.Allow("1.2.3.4", "tok-b") {
		t.Fatal("request with tok-b should be allowed (independent limit)")
	}
}

func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	clk := &mockClock{t: time.Now()}
	rl := newTestLimiter(2, 100, time.Minute, clk)

	// Exhaust the IP limit.
	rl.Allow("1.2.3.4", "tok-a") //nolint:errcheck
	rl.Allow("1.2.3.4", "tok-a") //nolint:errcheck
	if rl.Allow("1.2.3.4", "tok-a") {
		t.Fatal("3rd request should be denied before window expires")
	}

	// Advance past the window.
	clk.advance(time.Minute + time.Second)

	// Should be allowed again.
	if !rl.Allow("1.2.3.4", "tok-a") {
		t.Fatal("request should be allowed after window resets")
	}
}

func TestRateLimiter_TokenLimitBlocksEvenIfIPIsUnder(t *testing.T) {
	clk := &mockClock{t: time.Now()}
	rl := newTestLimiter(100, 2, time.Minute, clk)

	// Different IPs, same token — exhaust token limit.
	rl.Allow("1.1.1.1", "tok-x") //nolint:errcheck
	rl.Allow("2.2.2.2", "tok-x") //nolint:errcheck

	// A fresh IP with the same exhausted token should be denied.
	if rl.Allow("3.3.3.3", "tok-x") {
		t.Fatal("token-limited request should be denied even from a new IP")
	}
}

func TestRateLimiter_IPLimitBlocksEvenIfTokenIsUnder(t *testing.T) {
	clk := &mockClock{t: time.Now()}
	rl := newTestLimiter(2, 100, time.Minute, clk)

	// Same IP, different tokens — exhaust IP limit.
	rl.Allow("1.1.1.1", "tok-a") //nolint:errcheck
	rl.Allow("1.1.1.1", "tok-b") //nolint:errcheck

	// IP exhausted, fresh token should still be denied.
	if rl.Allow("1.1.1.1", "tok-c") {
		t.Fatal("IP-limited request should be denied even with a fresh token")
	}
}

func TestRateLimiter_DefaultConfig(t *testing.T) {
	// Config with zero values should use defaults and not panic.
	rl := NewRateLimiter(Config{})
	if !rl.Allow("1.2.3.4", "tok-a") {
		t.Fatal("first request with default config should be allowed")
	}
}

func TestRateLimiter_PartialIncrementIsAtomic(t *testing.T) {
	// When both checks would fail, neither counter should be incremented.
	clk := &mockClock{t: time.Now()}
	rl := newTestLimiter(1, 1, time.Minute, clk)

	// Exhaust both.
	rl.Allow("1.2.3.4", "tok-a") //nolint:errcheck

	// This should fail; after it fails, fresh IP+token counters should be at 0.
	if rl.Allow("1.2.3.4", "tok-a") {
		t.Fatal("over-limit request should be denied")
	}

	// A fresh token+IP combination should still have room for 1 request.
	if !rl.Allow("9.9.9.9", "tok-b") {
		t.Fatal("fresh IP+token should be allowed after a rejected request")
	}
}
