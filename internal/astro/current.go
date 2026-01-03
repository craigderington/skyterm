package astro

import "time"

// CurrentTime returns the current time
// This is a helper to avoid import cycles
func CurrentTime() time.Time {
	return time.Now()
}
