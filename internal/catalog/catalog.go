package catalog

// Star represents a celestial object in the catalog
type Star struct {
	Name         string
	Magnitude    float64
	RA           float64 // Right Ascension in hours (0-24)
	Dec          float64 // Declination in degrees (-90 to +90)
	Altitude     float64 // Calculated altitude for observer
	Azimuth      float64 // Calculated azimuth for observer
	SpectralType rune
}

// LoadDefaultStars returns a hardcoded set of bright stars
// Kept for backward compatibility
func LoadDefaultStars() []Star {
	catalog := NewStarCatalog()
	return catalog.Stars()
}
