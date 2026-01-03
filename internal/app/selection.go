package app

import (
	"math"

	"github.com/craigderington/skyterm/internal/astro"
	"github.com/craigderington/skyterm/internal/catalog"
	"github.com/craigderington/skyterm/internal/render"
)

// SelectedObject represents the currently selected celestial object
type SelectedObject struct {
	Type string // "star", "planet", "deepsky", "constellation"
	Name string

	// Object-specific data
	Star     *catalog.Star
	Planet   *astro.Planet
	DeepSky  *catalog.MessierObject
}

// ClearSelection clears the current selection
func (m *Model) ClearSelection() {
	m.selectedObject = nil
	m.following = false
}

// SelectNearestObject finds and selects the nearest object to the view center
func (m *Model) SelectNearestObject() {
	if m.canvas == nil {
		return
	}

	centerX := m.canvas.Width / 2
	centerY := m.canvas.Height / 2

	minDist := math.MaxFloat64
	var nearest *SelectedObject

	// Check stars
	for _, star := range m.starCatalog.Stars() {
		if star.Magnitude > m.magnitudeLimit {
			continue
		}

		dist := m.distanceToObject(star.Altitude, star.Azimuth, centerX, centerY)
		if dist < minDist && dist < 15.0 { // Within 15 pixel radius
			minDist = dist
			starCopy := star
			nearest = &SelectedObject{
				Type: "star",
				Name: star.Name,
				Star: &starCopy,
			}
		}
	}

	// Check planets (if visible)
	if m.showPlanets && m.planetarySystem != nil {
		for _, planet := range m.planetarySystem.AllPlanets() {
			dist := m.distanceToObject(planet.Altitude, planet.Azimuth, centerX, centerY)
			if dist < minDist && dist < 15.0 {
				minDist = dist
				planetCopy := planet
				nearest = &SelectedObject{
					Type:   "planet",
					Name:   planet.Name,
					Planet: &planetCopy,
				}
			}
		}
	}

	// Check deep sky (if visible)
	if m.showDeepSky {
		for _, obj := range m.deepSkyCatalog.Objects() {
			if obj.Magnitude > m.magnitudeLimit+3 {
				continue
			}

			dist := m.distanceToObject(obj.Altitude, obj.Azimuth, centerX, centerY)
			if dist < minDist && dist < 15.0 {
				minDist = dist
				objCopy := obj
				nearest = &SelectedObject{
					Type:    "deepsky",
					Name:    obj.Name,
					DeepSky: &objCopy,
				}
			}
		}
	}

	m.selectedObject = nearest
}

// CenterOnSelected centers the view on the selected object
func (m *Model) CenterOnSelected() {
	if m.selectedObject == nil {
		return
	}

	var alt, az float64

	switch m.selectedObject.Type {
	case "star":
		if m.selectedObject.Star != nil {
			alt = m.selectedObject.Star.Altitude
			az = m.selectedObject.Star.Azimuth
		}
	case "planet":
		if m.selectedObject.Planet != nil {
			alt = m.selectedObject.Planet.Altitude
			az = m.selectedObject.Planet.Azimuth
		}
	case "deepsky":
		if m.selectedObject.DeepSky != nil {
			alt = m.selectedObject.DeepSky.Altitude
			az = m.selectedObject.DeepSky.Azimuth
		}
	}

	if alt != 0 || az != 0 {
		m.altitude = alt
		m.azimuth = az
	}
}

// UpdateFollowing updates view to follow selected object
func (m *Model) UpdateFollowing() {
	if !m.following || m.selectedObject == nil {
		return
	}

	// Update selected object's position
	switch m.selectedObject.Type {
	case "star":
		for _, star := range m.starCatalog.Stars() {
			if star.Name == m.selectedObject.Name {
				starCopy := star
				m.selectedObject.Star = &starCopy
				m.altitude = star.Altitude
				m.azimuth = star.Azimuth
				break
			}
		}
	case "planet":
		for _, planet := range m.planetarySystem.AllPlanets() {
			if planet.Name == m.selectedObject.Name {
				planetCopy := planet
				m.selectedObject.Planet = &planetCopy
				m.altitude = planet.Altitude
				m.azimuth = planet.Azimuth
				break
			}
		}
	case "deepsky":
		for _, obj := range m.deepSkyCatalog.Objects() {
			if obj.Name == m.selectedObject.Name {
				objCopy := obj
				m.selectedObject.DeepSky = &objCopy
				m.altitude = obj.Altitude
				m.azimuth = obj.Azimuth
				break
			}
		}
	}
}

// distanceToObject calculates screen distance from screen center to object position
func (m *Model) distanceToObject(alt, az float64, centerX, centerY int) float64 {
	// Project object to screen coordinates
	x, y, visible := render.Project(alt, az, m.altitude, m.azimuth, m.fov, m.canvas.Width, m.canvas.Height)

	// If not visible, return max distance
	if !visible {
		return math.MaxFloat64
	}

	// Calculate pixel distance from center
	dx := float64(x - centerX)
	dy := float64(y - centerY)
	return math.Sqrt(dx*dx + dy*dy)
}
