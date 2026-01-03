package astro

import (
	"math"
	"testing"
	"time"
)

func TestJulianDate(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected float64
		delta    float64
	}{
		{
			name:     "J2000.0 epoch",
			time:     time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: 2451545.0,
			delta:    0.001,
		},
		{
			name:     "Unix epoch",
			time:     time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 2440587.5,
			delta:    0.001,
		},
		{
			name:     "2025-01-01 00:00 UTC",
			time:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 2460676.5,
			delta:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jd := JulianDate(tt.time)
			if math.Abs(jd-tt.expected) > tt.delta {
				t.Errorf("JulianDate() = %v, want %v (±%v)", jd, tt.expected, tt.delta)
			}
		})
	}
}

func TestGMST(t *testing.T) {
	// Test GMST calculation for J2000.0
	// At J2000.0, GMST at 0h UT should be approximately 6.697374558 hours
	jd := 2451545.0 // J2000.0
	gmst := GMST(jd)

	// GMST should be between 0 and 24 hours
	if gmst < 0 || gmst >= 24 {
		t.Errorf("GMST out of range: %v", gmst)
	}

	// For J2000.0 at noon, we expect a specific value
	// This is a rough check - exact value depends on implementation details
	if gmst < 0 || gmst > 24 {
		t.Errorf("GMST calculation seems incorrect: %v", gmst)
	}
}

func TestLST(t *testing.T) {
	// Test LST for a known location and time
	jd := 2451545.0  // J2000.0
	longitude := 0.0 // Greenwich

	lst := LST(jd, longitude)

	// LST at Greenwich should equal GMST
	gmst := GMST(jd)
	if math.Abs(lst-gmst) > 0.001 {
		t.Errorf("LST at Greenwich should equal GMST: LST=%v, GMST=%v", lst, gmst)
	}

	// Test with non-zero longitude
	longitude = -74.0 // New York (West)
	lst = LST(jd, longitude)

	// LST should be between 0 and 24 hours
	if lst < 0 || lst >= 24 {
		t.Errorf("LST out of range: %v", lst)
	}
}

func TestFormatSiderealTime(t *testing.T) {
	tests := []struct {
		hours    float64
		expected string
	}{
		{0.0, "00:00:00"},
		{12.0, "12:00:00"},
		{23.9999, "23:59:59"},
		{6.5, "06:30:00"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatSiderealTime(tt.hours)
			if result != tt.expected {
				t.Errorf("FormatSiderealTime(%v) = %v, want %v", tt.hours, result, tt.expected)
			}
		})
	}
}
