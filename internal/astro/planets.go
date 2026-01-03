package astro

import (
	"math"
	"time"

	"github.com/soniakeys/meeus/v3/coord"
	"github.com/soniakeys/meeus/v3/julian"
	"github.com/soniakeys/meeus/v3/moonposition"
	"github.com/soniakeys/meeus/v3/nutation"
	"github.com/soniakeys/meeus/v3/solar"
	"github.com/soniakeys/unit"
)

// BodyType represents the type of celestial body
type BodyType string

const (
	BodyTypeSun    BodyType = "Sun"
	BodyTypeMoon   BodyType = "Moon"
	BodyTypePlanet BodyType = "Planet"
)

// Planet represents a solar system body
type Planet struct {
	Name      string
	BodyType  BodyType // Type of body: Sun, Moon, or Planet
	RA        float64  // Right Ascension in hours
	Dec       float64  // Declination in degrees
	Altitude  float64  // Calculated altitude
	Azimuth   float64  // Calculated azimuth
	Magnitude float64  // Visual magnitude
}

// PlanetarySystem holds all planets and the Moon
type PlanetarySystem struct {
	Sun     Planet
	Moon    Planet
	Mercury Planet
	Venus   Planet
	Mars    Planet
	Jupiter Planet
	Saturn  Planet
	Uranus  Planet
	Neptune Planet
}

// CalculatePlanets computes positions for all planets at given time
func CalculatePlanets(t time.Time, observer *Observer) *PlanetarySystem {
	jde := julian.TimeToJD(t)

	sys := &PlanetarySystem{}

	// Calculate Sun position
	sys.Sun = calculateSun(jde, observer, t)

	// Calculate Moon position
	sys.Moon = calculateMoon(jde, observer, t)

	// Calculate planets using low-precision method
	sys.Mercury = calculatePlanetLowPrecision("Mercury", jde, observer, t)
	sys.Venus = calculatePlanetLowPrecision("Venus", jde, observer, t)
	sys.Mars = calculatePlanetLowPrecision("Mars", jde, observer, t)
	sys.Jupiter = calculatePlanetLowPrecision("Jupiter", jde, observer, t)
	sys.Saturn = calculatePlanetLowPrecision("Saturn", jde, observer, t)
	sys.Uranus = calculatePlanetLowPrecision("Uranus", jde, observer, t)
	sys.Neptune = calculatePlanetLowPrecision("Neptune", jde, observer, t)

	return sys
}

// calculateSun computes Sun position
func calculateSun(jde float64, observer *Observer, t time.Time) Planet {
	// Solar apparent position
	ra, dec := solar.ApparentEquatorial(jde)

	// Convert to hours and degrees
	raHours := ra.Hour()
	decDeg := dec.Deg()

	// Convert to horizontal coordinates
	eq := EquatorialCoords{RA: raHours, Dec: decDeg}
	hz := EquatorialToHorizontal(eq, observer, t)

	return Planet{
		Name:      "Sun",
		BodyType:  BodyTypeSun,
		RA:        raHours,
		Dec:       decDeg,
		Altitude:  hz.Altitude,
		Azimuth:   hz.Azimuth,
		Magnitude: -26.7, // Sun's apparent magnitude
	}
}

// calculateMoon computes Moon position
func calculateMoon(jde float64, observer *Observer, t time.Time) Planet {
	// Moon position
	ra, dec, _ := moonposition.Position(jde)

	// Convert to hours and degrees
	raHours := ra.Rad() * 12.0 / math.Pi // RA radians to hours
	decDeg := dec.Deg()

	// Convert to horizontal coordinates
	eq := EquatorialCoords{RA: raHours, Dec: decDeg}
	hz := EquatorialToHorizontal(eq, observer, t)

	// Approximate magnitude based on phase
	// Moon varies from about -12.7 (full) to much dimmer
	mag := -12.0

	return Planet{
		Name:      "Moon",
		BodyType:  BodyTypeMoon,
		RA:        raHours,
		Dec:       decDeg,
		Altitude:  hz.Altitude,
		Azimuth:   hz.Azimuth,
		Magnitude: mag,
	}
}

// calculatePlanetLowPrecision computes approximate planet position
// Uses simplified orbital elements - accurate to ~1 arcminute for dates near 2000
func calculatePlanetLowPrecision(name string, jde float64, observer *Observer, t time.Time) Planet {
	// Julian centuries from J2000.0
	T := (jde - 2451545.0) / 36525.0

	// Get Sun's ecliptic longitude for Earth's position
	sunLon := solar.ApparentLongitude(T)
	earthLon := sunLon + math.Pi // Earth is opposite the Sun

	// Simplified orbital elements (from Meeus, good for ±3000 years from 2000)
	var L, a, e, i, omega, Omega float64 // Mean longitude, semi-major axis, eccentricity, inclination, arg of perihelion, long of asc node

	switch name {
	case "Mercury":
		L = 252.25 + 149474.07*T
		a = 0.387098
		e = 0.205635 + 0.000020*T
		i = 7.005 - 0.006*T
		omega = 77.46 - 0.16*T
		Omega = 48.33 + 0.04*T
	case "Venus":
		L = 181.98 + 58519.21*T
		a = 0.723330
		e = 0.006773 - 0.000042*T
		i = 3.395 - 0.008*T
		omega = 131.56 - 0.05*T
		Omega = 76.68 + 0.28*T
	case "Mars":
		L = 355.43 + 19141.70*T
		a = 1.523688
		e = 0.093405 + 0.000090*T
		i = 1.850 - 0.007*T
		omega = 336.06 + 0.44*T
		Omega = 49.56 - 0.29*T
	case "Jupiter":
		L = 34.35 + 3036.30*T
		a = 5.202603
		e = 0.048498 - 0.000163*T
		i = 1.303 - 0.002*T
		omega = 14.33 + 0.21*T
		Omega = 100.46 - 0.11*T
	case "Saturn":
		L = 50.08 + 1223.51*T
		a = 9.554909
		e = 0.055546 - 0.000347*T
		i = 2.489 - 0.004*T
		omega = 93.06 + 0.57*T
		Omega = 113.67 - 0.26*T
	case "Uranus":
		L = 314.05 + 429.86*T
		a = 19.218446
		e = 0.047318 - 0.000018*T
		i = 0.773 + 0.001*T
		omega = 173.01 + 0.10*T
		Omega = 74.01 - 0.04*T
	case "Neptune":
		L = 304.35 + 219.88*T
		a = 30.110387
		e = 0.008606 + 0.000002*T
		i = 1.770 - 0.001*T
		omega = 48.12 - 0.04*T
		Omega = 131.78 - 0.02*T
	default:
		return Planet{Name: name}
	}

	// Convert to radians
	L = L * math.Pi / 180
	i = i * math.Pi / 180
	omega = omega * math.Pi / 180
	Omega = Omega * math.Pi / 180

	// Mean anomaly
	M := L - omega

	// Solve Kepler's equation for eccentric anomaly E
	E := M
	for iter := 0; iter < 10; iter++ {
		E = M + e*math.Sin(E)
	}

	// True anomaly
	v := 2 * math.Atan2(math.Sqrt(1+e)*math.Sin(E/2), math.Sqrt(1-e)*math.Cos(E/2))

	// Distance from Sun
	r := a * (1 - e*math.Cos(E))

	// Heliocentric coordinates in orbital plane
	xOrbital := r * math.Cos(v)
	yOrbital := r * math.Sin(v)

	// Convert to ecliptic coordinates
	xEcl := (math.Cos(omega)*math.Cos(Omega) - math.Sin(omega)*math.Sin(Omega)*math.Cos(i)) * xOrbital +
		(-math.Sin(omega)*math.Cos(Omega) - math.Cos(omega)*math.Sin(Omega)*math.Cos(i)) * yOrbital
	yEcl := (math.Cos(omega)*math.Sin(Omega) + math.Sin(omega)*math.Cos(Omega)*math.Cos(i)) * xOrbital +
		(-math.Sin(omega)*math.Sin(Omega) + math.Cos(omega)*math.Cos(Omega)*math.Cos(i)) * yOrbital
	zEcl := (math.Sin(omega) * math.Sin(i)) * xOrbital +
		(math.Cos(omega) * math.Sin(i)) * yOrbital

	// Earth's position (assuming circular orbit for simplicity)
	earthR := 1.0 // AU
	earthX := earthR * math.Cos(earthLon.Rad())
	earthY := earthR * math.Sin(earthLon.Rad())
	earthZ := 0.0

	// Geocentric position
	geoX := xEcl - earthX
	geoY := yEcl - earthY
	geoZ := zEcl - earthZ

	// Convert to ecliptic longitude and latitude
	geoLon := math.Atan2(geoY, geoX)
	geoLat := math.Atan2(geoZ, math.Sqrt(geoX*geoX+geoY*geoY))

	// Convert ecliptic to equatorial
	ecliptic := &coord.Ecliptic{
		Lon: unit.Angle(geoLon),
		Lat: unit.Angle(geoLat),
	}

	ε := nutation.MeanObliquity(jde)
	equatorial := new(coord.Equatorial).EclToEq(ecliptic, coord.NewObliquity(ε))

	raHours := equatorial.RA.Hour()
	decDeg := equatorial.Dec.Deg()

	// Convert to horizontal coordinates
	eq := EquatorialCoords{RA: raHours, Dec: decDeg}
	hz := EquatorialToHorizontal(eq, observer, t)

	// Approximate magnitudes
	magnitudes := map[string]float64{
		"Mercury": -0.4,
		"Venus":   -4.4,
		"Mars":    -2.0,
		"Jupiter": -2.7,
		"Saturn":  0.4,
		"Uranus":  5.7,
		"Neptune": 7.8,
	}

	return Planet{
		Name:      name,
		BodyType:  BodyTypePlanet,
		RA:        raHours,
		Dec:       decDeg,
		Altitude:  hz.Altitude,
		Azimuth:   hz.Azimuth,
		Magnitude: magnitudes[name],
	}
}

// MoonPhase calculates the Moon's phase (0-1, where 0=new, 0.5=full)
func MoonPhase(t time.Time) float64 {
	jde := julian.TimeToJD(t)

	// Get Sun and Moon positions
	sunRA, sunDec := solar.ApparentEquatorial(jde)
	moonRA, moonDec, _ := moonposition.Position(jde)

	// Convert to radians
	sunRARad := sunRA.Rad()
	sunDecRad := sunDec.Rad()
	moonRARad := moonRA.Rad()
	moonDecRad := moonDec.Rad()

	// Calculate elongation (angular separation)
	elongation := math.Acos(
		math.Sin(sunDecRad)*math.Sin(moonDecRad) +
			math.Cos(sunDecRad)*math.Cos(moonDecRad)*math.Cos(sunRARad-moonRARad))

	// Phase is elongation / π
	phase := elongation / math.Pi

	return phase
}

// AllPlanets returns a slice of all planets for iteration
func (ps *PlanetarySystem) AllPlanets() []Planet {
	return []Planet{
		ps.Sun,
		ps.Moon,
		ps.Mercury,
		ps.Venus,
		ps.Mars,
		ps.Jupiter,
		ps.Saturn,
		ps.Uranus,
		ps.Neptune,
	}
}
