package config

import (
	"os"
	"path/filepath"

	"github.com/craigderington/skyterm/internal/astro"
	"gopkg.in/yaml.v3"
)

// Config holds application configuration
type Config struct {
	Location LocationConfig
	Display  DisplayConfig
	Time     TimeConfig
	Controls ControlsConfig
}

// LocationConfig holds observer location settings
type LocationConfig struct {
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
	Altitude  float64 `yaml:"altitude"`
	Name      string  `yaml:"name"`
}

// DisplayConfig holds display settings
type DisplayConfig struct {
	MagnitudeLimit         float64 `yaml:"magnitude_limit"`
	ShowConstellationLines bool    `yaml:"show_constellation_lines"`
	ShowConstellationNames bool    `yaml:"show_constellation_names"`
	ShowCoordinateGrid     bool    `yaml:"show_coordinate_grid"`
	ShowPlanetLabels       bool    `yaml:"show_planet_labels"`
	ColorStarsByType       bool    `yaml:"color_stars_by_type"`
	UseBrailleRendering    bool    `yaml:"use_braille_rendering"`
}

// TimeConfig holds time-related settings
type TimeConfig struct {
	UseUTC   bool   `yaml:"use_utc"`
	TimeStep string `yaml:"time_step"`
}

// ControlsConfig holds control settings
type ControlsConfig struct {
	PanSpeed         float64 `yaml:"pan_speed"`
	FastPanMultiplier float64 `yaml:"fast_pan_multiplier"`
	ZoomStep         float64 `yaml:"zoom_step"`
}

// Load loads configuration from XDG config directory
// Falls back to defaults if config file doesn't exist
func Load() (*Config, error) {
	configPath := getConfigPath()

	// If config file doesn't exist, return defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), nil
	}

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), nil
	}

	// Merge with defaults for any missing values
	defaults := DefaultConfig()
	if cfg.Location.Latitude == 0 && cfg.Location.Longitude == 0 {
		cfg.Location = defaults.Location
	}
	if cfg.Display.MagnitudeLimit == 0 {
		cfg.Display = defaults.Display
	}
	if cfg.Controls.PanSpeed == 0 {
		cfg.Controls = defaults.Controls
	}

	return &cfg, nil
}

// Observer creates an Observer from the location configuration
func (c *Config) Observer() *astro.Observer {
	return astro.NewObserver(
		c.Location.Latitude,
		c.Location.Longitude,
		c.Location.Altitude,
		c.Location.Name,
	)
}

// getConfigPath returns the XDG config path for skyterm
func getConfigPath() string {
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		home, _ := os.UserHomeDir()
		xdgConfig = filepath.Join(home, ".config")
	}
	return filepath.Join(xdgConfig, "skyterm", "config.yaml")
}
