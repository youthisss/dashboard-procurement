package handlers

import (
	"strings"
	"sync"
	"time"
)

const (
	loginAttemptLimit    = 5
	loginLockout         = 15 * time.Minute
	loginAttemptStoreCap = 10000
	loginAttemptEntryTTL = 24 * time.Hour
)

type loginLimiter struct {
	mu       sync.Mutex
	now      func() time.Time
	attempts map[string]loginAttempt
}

type loginAttempt struct {
	Failures    int
	LockedUntil time.Time
	LastSeen    time.Time
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
	now := l.now()
	l.prune(now)

	entry, ok := l.attempts[key]
	if !ok {
		return false
	}
	if entry.LockedUntil.After(now) {
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
	now := l.now()
	l.prune(now)

	entry := l.attempts[key]
	entry.Failures++
	entry.LastSeen = now
	if entry.Failures >= loginAttemptLimit {
		entry.LockedUntil = now.Add(loginLockout)
	}
	l.attempts[key] = entry
	if overflow := len(l.attempts) - loginAttemptStoreCap; overflow > 0 {
		l.evictOldest(overflow)
	}
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

func (l *loginLimiter) prune(now time.Time) {
	for key, entry := range l.attempts {
		if !entry.LockedUntil.IsZero() && !entry.LockedUntil.After(now) {
			delete(l.attempts, key)
			continue
		}
		if entry.LastSeen.IsZero() {
			if entry.LockedUntil.IsZero() {
				delete(l.attempts, key)
			}
			continue
		}
		if now.Sub(entry.LastSeen) > loginAttemptEntryTTL {
			delete(l.attempts, key)
		}
	}
}

func (l *loginLimiter) evictOldest(count int) {
	for i := 0; i < count; i++ {
		oldestKey := ""
		var oldestSeen time.Time
		for key, entry := range l.attempts {
			seen := entry.LastSeen
			if seen.IsZero() {
				seen = entry.LockedUntil
			}
			if oldestKey == "" || seen.Before(oldestSeen) {
				oldestKey = key
				oldestSeen = seen
			}
		}
		if oldestKey == "" {
			return
		}
		delete(l.attempts, oldestKey)
	}
}
