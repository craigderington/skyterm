package catalog

import (
	"time"

	"github.com/craigderington/skyterm/internal/astro"
)

// StarCatalog holds the star catalog and provides methods to work with it
type StarCatalog struct {
	stars []Star
}

// NewStarCatalog creates a new star catalog
func NewStarCatalog() *StarCatalog {
	return &StarCatalog{
		stars: loadBrightStars(),
	}
}

// UpdatePositions updates all star positions for the given observer and time
func (sc *StarCatalog) UpdatePositions(observer *astro.Observer, t time.Time) {
	for i := range sc.stars {
		eq := astro.EquatorialCoords{
			RA:  sc.stars[i].RA,
			Dec: sc.stars[i].Dec,
		}
		hz := astro.EquatorialToHorizontal(eq, observer, t)
		sc.stars[i].Altitude = hz.Altitude
		sc.stars[i].Azimuth = hz.Azimuth
	}
}

// Stars returns all stars in the catalog
func (sc *StarCatalog) Stars() []Star {
	return sc.stars
}

// loadBrightStars loads a comprehensive set of bright stars
func loadBrightStars() []Star {
	return []Star{
		// Top 20 brightest stars
		{Name: "Sirius", Magnitude: -1.46, RA: 6.75, Dec: -16.72, SpectralType: 'A'},
		{Name: "Canopus", Magnitude: -0.74, RA: 6.40, Dec: -52.70, SpectralType: 'F'},
		{Name: "Arcturus", Magnitude: -0.05, RA: 14.26, Dec: 19.18, SpectralType: 'K'},
		{Name: "Rigel Kentaurus", Magnitude: -0.01, RA: 14.66, Dec: -60.83, SpectralType: 'G'},
		{Name: "Vega", Magnitude: 0.03, RA: 18.62, Dec: 38.78, SpectralType: 'A'},
		{Name: "Capella", Magnitude: 0.08, RA: 5.28, Dec: 46.00, SpectralType: 'G'},
		{Name: "Rigel", Magnitude: 0.13, RA: 5.24, Dec: -8.20, SpectralType: 'B'},
		{Name: "Procyon", Magnitude: 0.38, RA: 7.66, Dec: 5.22, SpectralType: 'F'},
		{Name: "Achernar", Magnitude: 0.45, RA: 1.63, Dec: -57.24, SpectralType: 'B'},
		{Name: "Betelgeuse", Magnitude: 0.50, RA: 5.92, Dec: 7.41, SpectralType: 'M'},
		{Name: "Hadar", Magnitude: 0.61, RA: 14.06, Dec: -60.37, SpectralType: 'B'},
		{Name: "Altair", Magnitude: 0.76, RA: 19.85, Dec: 8.87, SpectralType: 'A'},
		{Name: "Acrux", Magnitude: 0.77, RA: 12.44, Dec: -63.10, SpectralType: 'B'},
		{Name: "Aldebaran", Magnitude: 0.85, RA: 4.60, Dec: 16.51, SpectralType: 'K'},
		{Name: "Spica", Magnitude: 0.98, RA: 13.42, Dec: -11.16, SpectralType: 'B'},
		{Name: "Antares", Magnitude: 1.06, RA: 16.49, Dec: -26.43, SpectralType: 'M'},
		{Name: "Pollux", Magnitude: 1.14, RA: 7.75, Dec: 28.03, SpectralType: 'K'},
		{Name: "Fomalhaut", Magnitude: 1.16, RA: 22.96, Dec: -29.62, SpectralType: 'A'},
		{Name: "Deneb", Magnitude: 1.25, RA: 20.69, Dec: 45.28, SpectralType: 'A'},
		{Name: "Mimosa", Magnitude: 1.25, RA: 12.79, Dec: -59.69, SpectralType: 'B'},

		// Orion constellation
		{Name: "Bellatrix", Magnitude: 1.64, RA: 5.42, Dec: 6.35, SpectralType: 'B'},
		{Name: "Alnilam", Magnitude: 1.69, RA: 5.60, Dec: -1.20, SpectralType: 'B'},
		{Name: "Alnitak", Magnitude: 1.77, RA: 5.68, Dec: -1.94, SpectralType: 'O'},
		{Name: "Mintaka", Magnitude: 2.23, RA: 5.53, Dec: -0.30, SpectralType: 'B'},
		{Name: "Saiph", Magnitude: 2.06, RA: 5.80, Dec: -9.67, SpectralType: 'B'},

		// Big Dipper / Ursa Major
		{Name: "Alioth", Magnitude: 1.76, RA: 12.90, Dec: 55.96, SpectralType: 'A'},
		{Name: "Dubhe", Magnitude: 1.79, RA: 11.06, Dec: 61.75, SpectralType: 'K'},
		{Name: "Alkaid", Magnitude: 1.85, RA: 13.79, Dec: 49.31, SpectralType: 'B'},
		{Name: "Mizar", Magnitude: 2.23, RA: 13.40, Dec: 54.93, SpectralType: 'A'},
		{Name: "Merak", Magnitude: 2.34, RA: 11.03, Dec: 56.38, SpectralType: 'A'},
		{Name: "Phecda", Magnitude: 2.41, RA: 11.90, Dec: 53.69, SpectralType: 'A'},
		{Name: "Megrez", Magnitude: 3.32, RA: 12.26, Dec: 57.03, SpectralType: 'A'},

		// Leo
		{Name: "Regulus", Magnitude: 1.35, RA: 10.14, Dec: 11.97, SpectralType: 'B'},
		{Name: "Denebola", Magnitude: 2.14, RA: 11.82, Dec: 14.57, SpectralType: 'A'},
		{Name: "Algieba", Magnitude: 2.08, RA: 10.33, Dec: 19.84, SpectralType: 'K'},

		// Gemini
		{Name: "Castor", Magnitude: 1.58, RA: 7.58, Dec: 31.89, SpectralType: 'A'},
		{Name: "Alhena", Magnitude: 1.93, RA: 6.63, Dec: 16.40, SpectralType: 'A'},

		// Taurus
		{Name: "Elnath", Magnitude: 1.65, RA: 5.44, Dec: 28.61, SpectralType: 'B'},

		// Auriga
		{Name: "Menkalinan", Magnitude: 1.90, RA: 6.00, Dec: 44.95, SpectralType: 'A'},

		// Virgo
		{Name: "Porrima", Magnitude: 2.74, RA: 12.69, Dec: -1.45, SpectralType: 'F'},
		{Name: "Vindemiatrix", Magnitude: 2.85, RA: 13.04, Dec: 10.96, SpectralType: 'G'},

		// Bootes
		{Name: "Nekkar", Magnitude: 3.50, RA: 15.03, Dec: 40.39, SpectralType: 'G'},
		{Name: "Seginus", Magnitude: 3.04, RA: 14.53, Dec: 38.31, SpectralType: 'A'},

		// Scorpius
		{Name: "Shaula", Magnitude: 1.62, RA: 17.56, Dec: -37.10, SpectralType: 'B'},
		{Name: "Sargas", Magnitude: 1.86, RA: 17.62, Dec: -43.00, SpectralType: 'F'},
		{Name: "Dschubba", Magnitude: 2.32, RA: 16.00, Dec: -22.62, SpectralType: 'B'},

		// Sagittarius
		{Name: "Kaus Australis", Magnitude: 1.79, RA: 18.40, Dec: -34.38, SpectralType: 'B'},
		{Name: "Nunki", Magnitude: 2.05, RA: 18.92, Dec: -26.30, SpectralType: 'B'},

		// Aquila
		{Name: "Tarazed", Magnitude: 2.72, RA: 19.77, Dec: 10.61, SpectralType: 'K'},
		{Name: "Alshain", Magnitude: 3.71, RA: 19.92, Dec: 6.41, SpectralType: 'A'},

		// Cygnus
		{Name: "Sadr", Magnitude: 2.23, RA: 20.37, Dec: 40.26, SpectralType: 'F'},
		{Name: "Gienah", Magnitude: 2.48, RA: 20.77, Dec: 33.97, SpectralType: 'K'},

		// Lyra
		{Name: "Sheliak", Magnitude: 3.52, RA: 18.83, Dec: 33.36, SpectralType: 'B'},
		{Name: "Sulafat", Magnitude: 3.25, RA: 18.98, Dec: 32.69, SpectralType: 'B'},

		// Aquarius
		{Name: "Sadalsuud", Magnitude: 2.90, RA: 21.52, Dec: -5.57, SpectralType: 'G'},
		{Name: "Sadalmelik", Magnitude: 3.00, RA: 22.10, Dec: -0.32, SpectralType: 'G'},

		// Pegasus
		{Name: "Enif", Magnitude: 2.38, RA: 21.74, Dec: 9.88, SpectralType: 'K'},
		{Name: "Scheat", Magnitude: 2.44, RA: 23.06, Dec: 28.08, SpectralType: 'M'},
		{Name: "Markab", Magnitude: 2.49, RA: 23.08, Dec: 15.21, SpectralType: 'B'},
		{Name: "Algenib", Magnitude: 2.83, RA: 0.22, Dec: 15.18, SpectralType: 'B'},

		// Andromeda
		{Name: "Alpheratz", Magnitude: 2.06, RA: 0.14, Dec: 29.09, SpectralType: 'B'},
		{Name: "Mirach", Magnitude: 2.07, RA: 1.16, Dec: 35.62, SpectralType: 'M'},
		{Name: "Almach", Magnitude: 2.10, RA: 2.07, Dec: 42.33, SpectralType: 'K'},

		// Cassiopeia
		{Name: "Schedar", Magnitude: 2.24, RA: 0.67, Dec: 56.54, SpectralType: 'K'},
		{Name: "Caph", Magnitude: 2.28, RA: 0.15, Dec: 59.15, SpectralType: 'F'},
		{Name: "Ruchbah", Magnitude: 2.68, RA: 1.43, Dec: 60.24, SpectralType: 'A'},

		// Perseus
		{Name: "Mirfak", Magnitude: 1.79, RA: 3.41, Dec: 49.86, SpectralType: 'F'},
		{Name: "Algol", Magnitude: 2.09, RA: 3.14, Dec: 40.96, SpectralType: 'B'},

		// Aries
		{Name: "Hamal", Magnitude: 2.00, RA: 2.12, Dec: 23.46, SpectralType: 'K'},
		{Name: "Sheratan", Magnitude: 2.64, RA: 1.91, Dec: 20.81, SpectralType: 'A'},

		// Cetus
		{Name: "Deneb Kaitos", Magnitude: 2.04, RA: 0.73, Dec: -17.99, SpectralType: 'K'},
		{Name: "Menkar", Magnitude: 2.54, RA: 3.04, Dec: 4.09, SpectralType: 'M'},

		// Eridanus
		{Name: "Cursa", Magnitude: 2.79, RA: 5.13, Dec: -5.09, SpectralType: 'A'},
		{Name: "Zaurak", Magnitude: 3.03, RA: 3.97, Dec: -13.51, SpectralType: 'M'},

		// Canis Major
		{Name: "Adhara", Magnitude: 1.50, RA: 6.98, Dec: -28.97, SpectralType: 'B'},
		{Name: "Wezen", Magnitude: 1.83, RA: 7.14, Dec: -26.39, SpectralType: 'F'},
		{Name: "Mirzam", Magnitude: 1.98, RA: 6.38, Dec: -17.96, SpectralType: 'B'},
		{Name: "Aludra", Magnitude: 2.45, RA: 7.40, Dec: -29.30, SpectralType: 'B'},

		// Canis Minor - just Procyon already listed

		// Cancer
		{Name: "Acubens", Magnitude: 4.26, RA: 8.97, Dec: 11.86, SpectralType: 'A'},

		// Corvus
		{Name: "Gienah Corvi", Magnitude: 2.58, RA: 12.26, Dec: -17.54, SpectralType: 'B'},
		{Name: "Algorab", Magnitude: 2.94, RA: 12.50, Dec: -16.52, SpectralType: 'A'},

		// Centaurus
		{Name: "Menkent", Magnitude: 2.06, RA: 14.11, Dec: -36.37, SpectralType: 'K'},

		// Lupus
		{Name: "Men", Magnitude: 2.68, RA: 14.70, Dec: -43.13, SpectralType: 'B'},

		// Libra
		{Name: "Zubenelgenubi", Magnitude: 2.75, RA: 14.85, Dec: -16.04, SpectralType: 'A'},
		{Name: "Zubeneschamali", Magnitude: 2.61, RA: 15.28, Dec: -9.38, SpectralType: 'B'},

		// Ophiuchus
		{Name: "Rasalhague", Magnitude: 2.08, RA: 17.58, Dec: 12.56, SpectralType: 'A'},
		{Name: "Sabik", Magnitude: 2.43, RA: 17.17, Dec: -15.72, SpectralType: 'A'},

		// Hercules
		{Name: "Rasalgethi", Magnitude: 3.37, RA: 17.24, Dec: 14.39, SpectralType: 'M'},
		{Name: "Kornephoros", Magnitude: 2.78, RA: 16.50, Dec: 21.49, SpectralType: 'G'},

		// Corona Borealis
		{Name: "Alphecca", Magnitude: 2.22, RA: 15.58, Dec: 26.71, SpectralType: 'A'},

		// Southern Cross / Crux
		{Name: "Gacrux", Magnitude: 1.59, RA: 12.52, Dec: -57.11, SpectralType: 'M'},

		// Grus
		{Name: "Alnair", Magnitude: 1.73, RA: 22.14, Dec: -46.98, SpectralType: 'B'},
		{Name: "Al Dhanab", Magnitude: 3.01, RA: 22.71, Dec: -46.88, SpectralType: 'M'},

		// Pavo
		{Name: "Peacock", Magnitude: 1.94, RA: 20.43, Dec: -56.74, SpectralType: 'B'},

		// Tucana
		{Name: "Alpha Tucanae", Magnitude: 2.87, RA: 22.31, Dec: -60.26, SpectralType: 'K'},

		// Phoenix
		{Name: "Ankaa", Magnitude: 2.40, RA: 0.44, Dec: -42.31, SpectralType: 'K'},

		// Pisces Austrinus - Fomalhaut already listed

		// Polaris (North Star)
		{Name: "Polaris", Magnitude: 1.98, RA: 2.53, Dec: 89.26, SpectralType: 'F'},
	}
}
