package astro

import (
	"math"
	"time"
)

// EquatorialCoords represents equatorial coordinates (RA/Dec)
type EquatorialCoords struct {
	RA  float64 // Right Ascension in hours (0-24)
	Dec float64 // Declination in degrees (-90 to +90)
}

// HorizontalCoords represents horizontal coordinates (Alt/Az)
type HorizontalCoords struct {
	Altitude float64 // Altitude in degrees (-90 to +90)
	Azimuth  float64 // Azimuth in degrees (0-360, 0=North)
}

// EquatorialToHorizontal converts RA/Dec to Alt/Az for given observer and time
func EquatorialToHorizontal(eq EquatorialCoords, observer *Observer, t time.Time) HorizontalCoords {
	// Get Local Sidereal Time
	lst := observer.LST(t)

	// Convert to radians
	raRad := eq.RA * 15.0 * math.Pi / 180.0  // Hours to degrees to radians
	decRad := eq.Dec * math.Pi / 180.0
	latRad := observer.Latitude * math.Pi / 180.0
	lstRad := lst * 15.0 * math.Pi / 180.0

	// Calculate Hour Angle
	ha := lstRad - raRad

	// Calculate Altitude
	sinAlt := math.Sin(decRad)*math.Sin(latRad) +
	          math.Cos(decRad)*math.Cos(latRad)*math.Cos(ha)
	altitude := math.Asin(sinAlt)

	// Calculate Azimuth
	cosAz := (math.Sin(decRad) - math.Sin(latRad)*sinAlt) /
	         (math.Cos(latRad) * math.Cos(altitude))

	// Clamp to prevent NaN from floating point errors
	cosAz = math.Max(-1.0, math.Min(1.0, cosAz))
	azimuth := math.Acos(cosAz)

	// Adjust azimuth quadrant based on hour angle
	if math.Sin(ha) > 0 {
		azimuth = 2.0*math.Pi - azimuth
	}

	return HorizontalCoords{
		Altitude: altitude * 180.0 / math.Pi,
		Azimuth:  azimuth * 180.0 / math.Pi,
	}
}

// HorizontalToEquatorial converts Alt/Az to RA/Dec for given observer and time
func HorizontalToEquatorial(hz HorizontalCoords, observer *Observer, t time.Time) EquatorialCoords {
	// Get Local Sidereal Time
	lst := observer.LST(t)

	// Convert to radians
	altRad := hz.Altitude * math.Pi / 180.0
	azRad := hz.Azimuth * math.Pi / 180.0
	latRad := observer.Latitude * math.Pi / 180.0

	// Calculate Declination
	sinDec := math.Sin(altRad)*math.Sin(latRad) +
	          math.Cos(altRad)*math.Cos(latRad)*math.Cos(azRad)
	dec := math.Asin(sinDec)

	// Calculate Hour Angle
	cosHA := (math.Sin(altRad) - math.Sin(latRad)*sinDec) /
	         (math.Cos(latRad) * math.Cos(dec))
	cosHA = math.Max(-1.0, math.Min(1.0, cosHA))
	ha := math.Acos(cosHA)

	// Adjust hour angle quadrant based on azimuth
	if math.Sin(azRad) > 0 {
		ha = 2.0*math.Pi - ha
	}

	// Calculate RA from LST and HA
	lstRad := lst * 15.0 * math.Pi / 180.0
	ra := lstRad - ha

	// Normalize RA to 0-2π
	for ra < 0 {
		ra += 2.0 * math.Pi
	}
	for ra >= 2.0*math.Pi {
		ra -= 2.0 * math.Pi
	}

	return EquatorialCoords{
		RA:  (ra * 180.0 / math.Pi) / 15.0, // Convert to hours
		Dec: dec * 180.0 / math.Pi,
	}
}

// FormatRA formats Right Ascension in hours to HH:MM:SS.S format
func FormatRA(hours float64) string {
	h := int(hours)
	m := int((hours - float64(h)) * 60.0)
	s := ((hours - float64(h)) * 60.0 - float64(m)) * 60.0

	return formatInt(h, 2) + "h " + formatInt(m, 2) + "m " + formatFloat(s, 4, 1) + "s"
}

// FormatDec formats Declination in degrees to ±DD:MM:SS format
func FormatDec(degrees float64) string {
	sign := "+"
	if degrees < 0 {
		sign = "-"
		degrees = -degrees
	}

	d := int(degrees)
	m := int((degrees - float64(d)) * 60.0)
	s := ((degrees - float64(d)) * 60.0 - float64(m)) * 60.0

	return sign + formatInt(d, 2) + "° " + formatInt(m, 2) + "' " + formatFloat(s, 4, 1) + "\""
}

func formatFloat(f float64, width, precision int) string {
	// Simple float formatting
	i := int(f)
	frac := int((f - float64(i)) * math.Pow10(precision))
	return formatInt(i, width-precision-1) + "." + formatInt(frac, precision)
}
