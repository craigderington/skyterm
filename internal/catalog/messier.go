package catalog

// MessierObject represents a deep sky object from the Messier catalog
type MessierObject struct {
	Number      int
	Name        string
	CommonName  string
	Type        string // Galaxy, Nebula, Cluster, etc.
	RA          float64 // Right Ascension in hours
	Dec         float64 // Declination in degrees
	Magnitude   float64
	Altitude    float64 // Calculated
	Azimuth     float64 // Calculated
}

// GetMessierCatalog returns notable Messier objects
func GetMessierCatalog() []MessierObject {
	return []MessierObject{
		// Most famous Messier objects
		{
			Number:     1,
			Name:       "M1",
			CommonName: "Crab Nebula",
			Type:       "Supernova Remnant",
			RA:         5.58,
			Dec:        22.02,
			Magnitude:  8.4,
		},
		{
			Number:     8,
			Name:       "M8",
			CommonName: "Lagoon Nebula",
			Type:       "Nebula",
			RA:         18.06,
			Dec:        -24.38,
			Magnitude:  6.0,
		},
		{
			Number:     13,
			Name:       "M13",
			CommonName: "Hercules Cluster",
			Type:       "Globular Cluster",
			RA:         16.69,
			Dec:        36.46,
			Magnitude:  5.8,
		},
		{
			Number:     31,
			Name:       "M31",
			CommonName: "Andromeda Galaxy",
			Type:       "Galaxy",
			RA:         0.71,
			Dec:        41.27,
			Magnitude:  3.4,
		},
		{
			Number:     42,
			Name:       "M42",
			CommonName: "Orion Nebula",
			Type:       "Nebula",
			RA:         5.59,
			Dec:        -5.39,
			Magnitude:  4.0,
		},
		{
			Number:     44,
			Name:       "M44",
			CommonName: "Beehive Cluster",
			Type:       "Open Cluster",
			RA:         8.67,
			Dec:        19.98,
			Magnitude:  3.7,
		},
		{
			Number:     45,
			Name:       "M45",
			CommonName: "Pleiades",
			Type:       "Open Cluster",
			RA:         3.79,
			Dec:        24.12,
			Magnitude:  1.6,
		},
		{
			Number:     51,
			Name:       "M51",
			CommonName: "Whirlpool Galaxy",
			Type:       "Galaxy",
			RA:         13.50,
			Dec:        47.20,
			Magnitude:  8.4,
		},
		{
			Number:     57,
			Name:       "M57",
			CommonName: "Ring Nebula",
			Type:       "Planetary Nebula",
			RA:         18.89,
			Dec:        33.03,
			Magnitude:  8.8,
		},
		{
			Number:     81,
			Name:       "M81",
			CommonName: "Bode's Galaxy",
			Type:       "Galaxy",
			RA:         9.93,
			Dec:        69.07,
			Magnitude:  6.9,
		},
	}
}
