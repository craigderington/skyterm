package astro

import (
	"math"
	"time"
)

// RiseSetTransit holds rise, set, and transit times for an object
type RiseSetTransit struct {
	Rise    *time.Time // nil if circumpolar (always up)
	Set     *time.Time // nil if circumpolar (always up)
	Transit time.Time  // When object crosses meridian (highest point)

	NeverRises bool // True if object never rises (always below horizon)
	Circumpolar bool // True if object never sets (always above horizon)
}

// CalculateRiseSetTransit calculates rise, set, and transit times for an object
// Uses simplified algorithm - good approximation for stars and deep sky objects
func CalculateRiseSetTransit(ra, dec float64, observer *Observer, t time.Time) RiseSetTransit {
	// Calculate for the given date
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	// Observer latitude in radians
	latRad := observer.Latitude * math.Pi / 180.0
	decRad := dec * math.Pi / 180.0

	// Check if object is circumpolar or never rises
	// cos(H0) = (sin(h0) - sin(φ)sin(δ)) / (cos(φ)cos(δ))
	// where h0 = -0.833° for stars (atmospheric refraction)
	h0 := -0.833 * math.Pi / 180.0 // Standard refraction

	cosH0 := (math.Sin(h0) - math.Sin(latRad)*math.Sin(decRad)) / (math.Cos(latRad) * math.Cos(decRad))

	if cosH0 > 1.0 {
		// Never rises
		return RiseSetTransit{
			NeverRises: true,
			Transit:    date,
		}
	}

	if cosH0 < -1.0 {
		// Circumpolar (always above horizon)
		transitTime := calculateTransit(ra, observer, date)
		return RiseSetTransit{
			Circumpolar: true,
			Transit:     transitTime,
		}
	}

	// Hour angle at rise/set
	H0 := math.Acos(cosH0) * 180.0 / math.Pi / 15.0 // Convert to hours

	// Transit time (when object crosses meridian)
	transitTime := calculateTransit(ra, observer, date)

	// Rise time = Transit - H0
	// Set time = Transit + H0
	riseTime := transitTime.Add(-time.Duration(H0 * float64(time.Hour)))
	setTime := transitTime.Add(time.Duration(H0 * float64(time.Hour)))

	return RiseSetTransit{
		Rise:    &riseTime,
		Set:     &setTime,
		Transit: transitTime,
	}
}

// calculateTransit calculates when object transits (crosses meridian)
func calculateTransit(ra float64, observer *Observer, date time.Time) time.Time {
	// Local sidereal time when object transits equals its RA
	// Convert RA to time of day

	// Get GMST at 0h UT for this date
	jd := JulianDate(date)
	gmst0 := GMST(jd)

	// Local sidereal time = GMST + longitude (in hours)
	lonHours := observer.Longitude / 15.0

	// Transit occurs when LST = RA
	// Time in hours from midnight = (RA - GMST0 - lonHours) * 0.997...
	// (sidereal day is ~4 minutes shorter than solar day)

	hoursSinceMidnight := ra - gmst0 - lonHours

	// Normalize to 0-24 range
	for hoursSinceMidnight < 0 {
		hoursSinceMidnight += 24.0
	}
	for hoursSinceMidnight >= 24.0 {
		hoursSinceMidnight -= 24.0
	}

	// Convert sidereal time to solar time (approximately)
	solarHours := hoursSinceMidnight * 0.99726957 // Sidereal to solar ratio

	hours := int(solarHours)
	minutes := int((solarHours - float64(hours)) * 60.0)

	transitTime := time.Date(
		date.Year(), date.Month(), date.Day(),
		hours, minutes, 0, 0,
		date.Location(),
	)

	return transitTime
}

// FormatTime formats a time pointer to HH:MM or a message if nil
func FormatTime(t *time.Time) string {
	if t == nil {
		return "---"
	}
	return t.Format("15:04")
}
