package app

import (
	"strings"
)

// performSearch searches for objects matching the query and selects the first match
func (m *Model) performSearch() {
	if m.searchQuery == "" {
		return
	}

	query := strings.ToLower(strings.TrimSpace(m.searchQuery))

	// Search stars
	for _, star := range m.starCatalog.Stars() {
		if strings.Contains(strings.ToLower(star.Name), query) {
			starCopy := star
			m.selectedObject = &SelectedObject{
				Type: "star",
				Name: star.Name,
				Star: &starCopy,
			}
			m.CenterOnSelected()
			m.showInfo = true
			return
		}
	}

	// Search planets
	if m.planetarySystem != nil {
		for _, planet := range m.planetarySystem.AllPlanets() {
			if strings.Contains(strings.ToLower(planet.Name), query) {
				planetCopy := planet
				m.selectedObject = &SelectedObject{
					Type:   "planet",
					Name:   planet.Name,
					Planet: &planetCopy,
				}
				m.CenterOnSelected()
				m.showInfo = true
				return
			}
		}
	}

	// Search deep sky
	for _, obj := range m.deepSkyCatalog.Objects() {
		objName := strings.ToLower(obj.Name + " " + obj.CommonName)
		if strings.Contains(objName, query) {
			objCopy := obj
			m.selectedObject = &SelectedObject{
				Type:    "deepsky",
				Name:    obj.Name,
				DeepSky: &objCopy,
			}
			m.CenterOnSelected()
			m.showInfo = true
			return
		}
	}
}
