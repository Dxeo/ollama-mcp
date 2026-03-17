package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math/rand/v2"
	"net"
	"strings"
	"syscall"
	"time"
)

const (
	retryMaxAttempts  = 4
	retryInitialDelay = 500 * time.Millisecond
	retryMaxDelay     = 10 * time.Second
	retryJitterFactor = 0.25
)

// withRetry executes fn with exponential backoff, retrying only on transient
// network errors (connection refused, timeout, EOF). Application-level errors
// are returned immediately without retry.
func withRetry[T any](ctx context.Context, op string, fn func() (T, error)) (T, error) {
	var zero T
	var lastErr error
	delay := retryInitialDelay

	for attempt := 1; attempt <= retryMaxAttempts; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		if !isRetryable(err) {
			return zero, err
		}

		lastErr = err

		if attempt == retryMaxAttempts {
			break
		}

		// Apply jitter: delay ± 25%
		jitter := delay / 4
		actual := delay - jitter + time.Duration(rand.Int64N(int64(2*jitter+1)))

		slog.Warn("retrying ollama call",
			"op", op,
			"attempt", attempt,
			"next_attempt", attempt+1,
			"backoff", actual.String(),
			"error", lastErr,
		)

		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-time.After(actual):
		}

		// Exponential backoff, capped at max
		delay *= 2
		if delay > retryMaxDelay {
			delay = retryMaxDelay
		}
	}

	return zero, lastErr
}

// isRetryable returns true for transient network errors that warrant a retry.
func isRetryable(err error) bool {
	// Connection refused
	if errors.Is(err, syscall.ECONNREFUSED) {
		return true
	}
	// Unexpected EOF (server dropped connection)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	// Network timeout
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	// Connection reset
	if errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	// Fallback: check error string for common transient patterns
	msg := err.Error()
	if strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "no such host") {
		return true
	}
	return false
}
