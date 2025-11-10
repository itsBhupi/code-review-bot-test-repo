package utils

import (
	"sync"
	"time"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	mu       sync.RWMutex
	attempts map[string]*attemptRecord
	
	// Configuration
	maxAttempts int
	window      time.Duration
	blockTime   time.Duration
}

type attemptRecord struct {
	count      int
	windowStart time.Time
	blockedUntil time.Time
}

// NewRateLimiter creates a new rate limiter
// maxAttempts: maximum login attempts allowed within window
// window: time window for counting attempts
// blockTime: how long to block after exceeding max attempts
func NewRateLimiter(maxAttempts int, window, blockTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		attempts:    make(map[string]*attemptRecord),
		maxAttempts: maxAttempts,
		window:      window,
		blockTime:   blockTime,
	}
	
	// Start cleanup goroutine
	go rl.cleanup()
	
	return rl
}

// Allow checks if a request from the given identifier is allowed
func (rl *RateLimiter) Allow(identifier string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	record, exists := rl.attempts[identifier]
	
	// If record doesn't exist, create it
	if !exists {
		rl.attempts[identifier] = &attemptRecord{
			count:       1,
			windowStart: now,
			blockedUntil: time.Time{},
		}
		return true
	}
	
	// Check if currently blocked
	if !record.blockedUntil.IsZero() && now.Before(record.blockedUntil) {
		return false
	}
	
	// Reset block if block time has passed
	if !record.blockedUntil.IsZero() && now.After(record.blockedUntil) {
		record.blockedUntil = time.Time{}
		record.count = 0
		record.windowStart = now
	}
	
	// Check if we need to reset the window
	if now.Sub(record.windowStart) > rl.window {
		record.count = 1
		record.windowStart = now
		return true
	}
	
	// Increment attempt count
	record.count++
	
	// Check if we've exceeded max attempts
	if record.count > rl.maxAttempts {
		record.blockedUntil = now.Add(rl.blockTime)
		return false
	}
	
	return true
}

// Reset clears the rate limit for a given identifier
func (rl *RateLimiter) Reset(identifier string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.attempts, identifier)
}

// cleanup periodically removes old records to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		
		for identifier, record := range rl.attempts {
			// Remove records that are old and not blocked
			if record.blockedUntil.IsZero() && now.Sub(record.windowStart) > rl.window*2 {
				delete(rl.attempts, identifier)
			}
			// Remove records where block time has expired
			if !record.blockedUntil.IsZero() && now.After(record.blockedUntil.Add(rl.window)) {
				delete(rl.attempts, identifier)
			}
		}
		
		rl.mu.Unlock()
	}
}

