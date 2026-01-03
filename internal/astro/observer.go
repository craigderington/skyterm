package astro

import "time"

// Observer represents an observer's location on Earth
type Observer struct {
	Latitude  float64 // Degrees, positive North
	Longitude float64 // Degrees, positive East
	Altitude  float64 // Meters above sea level
	Name      string  // Location name
}

// NewObserver creates a new observer at the given location
func NewObserver(lat, lon, alt float64, name string) *Observer {
	return &Observer{
		Latitude:  lat,
		Longitude: lon,
		Altitude:  alt,
		Name:      name,
	}
}

// LST returns the Local Sidereal Time for this observer at the given time
func (o *Observer) LST(t time.Time) float64 {
	jd := JulianDate(t)
	return LST(jd, o.Longitude)
}

// DefaultObserver returns a default observer location (New York City)
func DefaultObserver() *Observer {
	return NewObserver(40.7128, -74.0060, 10.0, "New York City")
}
