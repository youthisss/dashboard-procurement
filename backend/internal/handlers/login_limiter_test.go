package handlers

import (
	"fmt"
	"testing"
	"time"
)

func TestLoginLimiterLocksAfterRepeatedFailures(t *testing.T) {
	now := time.Date(2026, 4, 30, 10, 0, 0, 0, time.UTC)
	limiter := newLoginLimiter()
	limiter.now = func() time.Time { return now }

	for i := 0; i < loginAttemptLimit-1; i++ {
		limiter.recordFailure("192.0.2.10", "Admin")
		if limiter.locked("192.0.2.10", "admin") {
			t.Fatalf("limiter locked before attempt limit")
		}
	}

	limiter.recordFailure("192.0.2.10", "admin")
	if !limiter.locked("192.0.2.10", "ADMIN") {
		t.Fatalf("limiter did not lock after attempt limit")
	}
}

func TestLoginLimiterResetAfterSuccessfulLogin(t *testing.T) {
	now := time.Date(2026, 4, 30, 10, 0, 0, 0, time.UTC)
	limiter := newLoginLimiter()
	limiter.now = func() time.Time { return now }

	for i := 0; i < loginAttemptLimit; i++ {
		limiter.recordFailure("192.0.2.10", "admin")
	}
	if !limiter.locked("192.0.2.10", "admin") {
		t.Fatalf("limiter should be locked before reset")
	}

	limiter.reset("192.0.2.10", "admin")
	if limiter.locked("192.0.2.10", "admin") {
		t.Fatalf("limiter stayed locked after reset")
	}
}

func TestLoginLimiterCapsAttemptStoreSize(t *testing.T) {
	now := time.Date(2026, 4, 30, 10, 0, 0, 0, time.UTC)
	limiter := newLoginLimiter()
	limiter.now = func() time.Time { return now }

	for i := 0; i < loginAttemptStoreCap+50; i++ {
		limiter.recordFailure(fmt.Sprintf("192.0.2.%d", i), "admin")
	}

	if got := len(limiter.attempts); got > loginAttemptStoreCap {
		t.Fatalf("attempt store size = %d, want <= %d", got, loginAttemptStoreCap)
	}
}
