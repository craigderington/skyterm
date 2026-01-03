package astro

import (
	"math"
	"time"
)

// JulianDate converts a time.Time to Julian Date
func JulianDate(t time.Time) float64 {
	// Convert to UTC
	t = t.UTC()

	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()

	// Adjust for January and February
	if month <= 2 {
		year--
		month += 12
	}

	// Calculate day fraction
	dayFraction := float64(day) + float64(hour)/24.0 + float64(minute)/1440.0 + float64(second)/86400.0

	a := year / 100
	b := 2 - a + a/4

	jd := math.Floor(365.25*float64(year+4716)) + math.Floor(30.6001*float64(month+1)) + dayFraction + float64(b) - 1524.5

	return jd
}

// GMST calculates Greenwich Mean Sidereal Time in hours
// From "Astronomical Algorithms" by Jean Meeus
func GMST(jd float64) float64 {
	// Days since J2000.0
	d := jd - 2451545.0
	t := d / 36525.0

	// GMST at 0h UT
	gmst := 280.46061837 + 360.98564736629*d + 0.000387933*t*t - t*t*t/38710000.0

	// Normalize to 0-360 degrees
	gmst = math.Mod(gmst, 360.0)
	if gmst < 0 {
		gmst += 360.0
	}

	// Convert to hours
	return gmst / 15.0
}

// LST calculates Local Sidereal Time in hours
func LST(jd float64, longitude float64) float64 {
	gmst := GMST(jd)

	// Longitude in hours (positive East)
	lonHours := longitude / 15.0

	lst := gmst + lonHours

	// Normalize to 0-24 hours
	for lst < 0 {
		lst += 24.0
	}
	for lst >= 24.0 {
		lst -= 24.0
	}

	return lst
}

// DaysSinceJ2000 returns the number of days since J2000.0 epoch
func DaysSinceJ2000(t time.Time) float64 {
	jd := JulianDate(t)
	return jd - 2451545.0
}

// FormatSiderealTime formats sidereal time in hours to HH:MM:SS format
func FormatSiderealTime(hours float64) string {
	h := int(hours)
	m := int((hours - float64(h)) * 60.0)
	s := int(((hours - float64(h)) * 60.0 - float64(m)) * 60.0)
	return formatTime(h, m, s)
}

func formatTime(h, m, s int) string {
	return formatInt(h, 2) + ":" + formatInt(m, 2) + ":" + formatInt(s, 2)
}

func formatInt(n, width int) string {
	s := ""
	for i := 0; i < width; i++ {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
