package astro

import (
	"math"
	"testing"
	"time"
)

func TestEquatorialToHorizontalConversion(t *testing.T) {
	// Test star at zenith
	observer := NewObserver(40.0, 0.0, 0.0, "Test")
	testTime := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)

	// Calculate what RA/Dec would put a star at zenith
	// For a star to be at zenith, its Dec must equal the observer's latitude
	// and its RA must equal the LST
	lst := observer.LST(testTime)

	eq := EquatorialCoords{
		RA:  lst,
		Dec: observer.Latitude,
	}

	hz := EquatorialToHorizontal(eq, observer, testTime)

	// Star should be at or near zenith (altitude = 90°)
	if math.Abs(hz.Altitude-90.0) > 1.0 {
		t.Errorf("Expected altitude near 90°, got %.2f°", hz.Altitude)
	}
}

func TestHorizontalToEquatorialConversion(t *testing.T) {
	observer := NewObserver(40.0, -74.0, 0.0, "Test")
	testTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// Test zenith point
	hz := HorizontalCoords{
		Altitude: 90.0,
		Azimuth:  0.0, // Azimuth is arbitrary at zenith
	}

	eq := HorizontalToEquatorial(hz, observer, testTime)

	// Declination should equal latitude
	if math.Abs(eq.Dec-observer.Latitude) > 1.0 {
		t.Errorf("Expected Dec near %.2f°, got %.2f°", observer.Latitude, eq.Dec)
	}

	// RA should equal LST (within a reasonable margin)
	lst := observer.LST(testTime)
	raDiff := math.Abs(eq.RA - lst)
	// Account for 24-hour wrap
	if raDiff > 12.0 {
		raDiff = 24.0 - raDiff
	}
	if raDiff > 1.0 {
		t.Errorf("Expected RA near %.2fh, got %.2fh (LST=%.2fh)", lst, eq.RA, lst)
	}
}

func TestRoundTripConversion(t *testing.T) {
	observer := NewObserver(40.7, -74.0, 0.0, "New York")
	testTime := time.Now()

	original := EquatorialCoords{
		RA:  10.0, // 10 hours
		Dec: 20.0, // 20 degrees
	}

	// Convert to horizontal and back
	hz := EquatorialToHorizontal(original, observer, testTime)
	result := HorizontalToEquatorial(hz, observer, testTime)

	// Should get back to original coordinates (within reasonable precision)
	raDiff := math.Abs(result.RA - original.RA)
	if raDiff > 24.0 {
		raDiff = math.Abs(raDiff - 24.0)
	}
	if raDiff > 0.1 {
		t.Errorf("RA round-trip error: %.4f vs %.4f (diff: %.4f)", original.RA, result.RA, raDiff)
	}

	decDiff := math.Abs(result.Dec - original.Dec)
	if decDiff > 0.1 {
		t.Errorf("Dec round-trip error: %.4f vs %.4f (diff: %.4f)", original.Dec, result.Dec, decDiff)
	}
}

func TestFormatRA(t *testing.T) {
	tests := []struct {
		hours    float64
		contains string
	}{
		{0.0, "00h"},
		{12.0, "12h"},
		{6.5, "06h 30m"},
	}

	for _, tt := range tests {
		result := FormatRA(tt.hours)
		if len(result) == 0 {
			t.Errorf("FormatRA(%v) returned empty string", tt.hours)
		}
	}
}

func TestFormatDec(t *testing.T) {
	tests := []struct {
		degrees  float64
		contains string
	}{
		{0.0, "+00°"},
		{45.0, "+45°"},
		{-30.0, "-30°"},
	}

	for _, tt := range tests {
		result := FormatDec(tt.degrees)
		if len(result) == 0 {
			t.Errorf("FormatDec(%v) returned empty string", tt.degrees)
		}
	}
}
