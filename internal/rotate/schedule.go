package rotate

import (
	"fmt"
	"time"
)

// Schedule defines when automatic rotation should occur.
type Schedule struct {
	// Interval is the duration between rotations.
	Interval time.Duration
	// LastRotated is the timestamp of the most recent successful rotation.
	LastRotated time.Time
}

// IsDue reports whether a rotation is due based on the schedule.
func (s *Schedule) IsDue() bool {
	if s.Interval <= 0 {
		return false
	}
	if s.LastRotated.IsZero() {
		return true
	}
	return time.Since(s.LastRotated) >= s.Interval
}

// NextRotation returns the time at which the next rotation is due.
func (s *Schedule) NextRotation() time.Time {
	if s.LastRotated.IsZero() {
		return time.Now().UTC()
	}
	return s.LastRotated.Add(s.Interval)
}

// ParseInterval parses a human-readable interval string (e.g. "24h", "30m").
func ParseInterval(raw string) (time.Duration, error) {
	if raw == "" {
		return 0, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("rotate: invalid interval %q: %w", raw, err)
	}
	if d < 0 {
		return 0, fmt.Errorf("rotate: interval must be positive, got %q", raw)
	}
	return d, nil
}
