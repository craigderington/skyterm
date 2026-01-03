package config

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Location: LocationConfig{
			Latitude:  40.7128,
			Longitude: -74.0060,
			Altitude:  10.0,
			Name:      "New York City",
		},
		Display: DisplayConfig{
			MagnitudeLimit:         5.0,
			ShowConstellationLines: false,
			ShowConstellationNames: false,
			ShowCoordinateGrid:     false,
			ShowPlanetLabels:       false,
			ColorStarsByType:       true,
			UseBrailleRendering:    false,
		},
		Time: TimeConfig{
			UseUTC:   false,
			TimeStep: "1m",
		},
		Controls: ControlsConfig{
			PanSpeed:          5.0,
			FastPanMultiplier: 4.0,
			ZoomStep:          1.2,
		},
	}
}
