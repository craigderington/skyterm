package catalog

import (
	"time"

	"github.com/craigderington/skyterm/internal/astro"
)

// DeepSkyCatalog holds deep sky objects and provides methods to work with them
type DeepSkyCatalog struct {
	objects []MessierObject
}

// NewDeepSkyCatalog creates a new deep sky catalog
func NewDeepSkyCatalog() *DeepSkyCatalog {
	return &DeepSkyCatalog{
		objects: GetMessierCatalog(),
	}
}

// UpdatePositions updates all object positions for the given observer and time
func (dsc *DeepSkyCatalog) UpdatePositions(observer *astro.Observer, t time.Time) {
	for i := range dsc.objects {
		eq := astro.EquatorialCoords{
			RA:  dsc.objects[i].RA,
			Dec: dsc.objects[i].Dec,
		}
		hz := astro.EquatorialToHorizontal(eq, observer, t)
		dsc.objects[i].Altitude = hz.Altitude
		dsc.objects[i].Azimuth = hz.Azimuth
	}
}

// Objects returns all objects in the catalog
func (dsc *DeepSkyCatalog) Objects() []MessierObject {
	return dsc.objects
}
