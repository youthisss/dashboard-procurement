package handlers

import (
	"strings"
	"sync"
	"time"
)

const (
	loginAttemptLimit = 5
	loginLockout      = 15 * time.Minute
)

type loginLimiter struct {
	mu       sync.Mutex
	now      func() time.Time
	attempts map[string]loginAttempt
}

type loginAttempt struct {
	Failures    int
	LockedUntil time.Time
}

func newLoginLimiter() *loginLimiter {
	return &loginLimiter{
		now:      time.Now,
		attempts: make(map[string]loginAttempt),
	}
}

func (l *loginLimiter) locked(ip, username string) bool {
	if l == nil {
		return false
	}
	key := loginKey(ip, username)
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.attempts[key]
	if !ok {
		return false
	}
	if entry.LockedUntil.After(l.now()) {
		return true
	}
	if !entry.LockedUntil.IsZero() {
		delete(l.attempts, key)
	}
	return false
}

func (l *loginLimiter) recordFailure(ip, username string) {
	if l == nil {
		return
	}
	key := loginKey(ip, username)
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := l.attempts[key]
	entry.Failures++
	if entry.Failures >= loginAttemptLimit {
		entry.LockedUntil = l.now().Add(loginLockout)
	}
	l.attempts[key] = entry
}

func (l *loginLimiter) reset(ip, username string) {
	if l == nil {
		return
	}
	key := loginKey(ip, username)
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, key)
}

func loginKey(ip, username string) string {
	normalizedIP := strings.TrimSpace(strings.ToLower(ip))
	normalizedUsername := strings.TrimSpace(strings.ToLower(username))
	return normalizedIP + "|" + normalizedUsername
}
