package render

import (
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/craigderington/skyterm/internal/catalog"
)

// Star magnitude to Unicode character mapping
var magnitudeChars = map[int]rune{
	-2: '★', // Very bright (Sirius, Canopus)
	-1: '✦', // Bright
	0:  '●', // First magnitude
	1:  '◉', // Second magnitude
	2:  '○', // Third magnitude
	3:  '◦', // Fourth magnitude
	4:  '·', // Fifth magnitude
	5:  '⋅', // Sixth magnitude (limit of naked eye)
	6:  '·', // Dim
}

// Spectral type to color mapping
var spectralColors = map[rune]lipgloss.Color{
	'O': lipgloss.Color("27"),  // Hot blue
	'B': lipgloss.Color("75"),  // Blue-white
	'A': lipgloss.Color("231"), // White
	'F': lipgloss.Color("230"), // Yellow-white
	'G': lipgloss.Color("229"), // Yellow (Sun-like)
	'K': lipgloss.Color("214"), // Orange
	'M': lipgloss.Color("196"), // Cool red
}

func RenderStars(canvas *Canvas, stars []catalog.Star, centerAlt, centerAz, fov, magLimit float64) {
	for _, star := range stars {
		if star.Magnitude > magLimit {
			continue
		}

		// Project star to screen coordinates
		x, y, visible := Project(star.Altitude, star.Azimuth, centerAlt, centerAz, fov, canvas.Width, canvas.Height)
		if !visible {
			continue
		}

		// Get character for magnitude
		char := getCharForMagnitude(star.Magnitude)

		// Get enhanced style (color + bold for bright stars)
		style := GetStarStyle(star.Magnitude, star.SpectralType, true)

		canvas.Set(x, y, char, style)
	}
}

func getCharForMagnitude(mag float64) rune {
	// Round magnitude to nearest integer
	magInt := int(math.Round(mag))
	if char, ok := magnitudeChars[magInt]; ok {
		return char
	}
	return '·'
}

func getColorForSpectralType(spectralType rune) lipgloss.Color {
	if color, ok := spectralColors[spectralType]; ok {
		return color
	}
	return lipgloss.Color("231") // Default to white
}

// Project performs stereographic projection from celestial coordinates to screen coordinates
func Project(alt, az, centerAlt, centerAz, fov float64, screenWidth, screenHeight int) (x, y int, visible bool) {
	// Convert degrees to radians
	altRad := alt * math.Pi / 180.0
	azRad := az * math.Pi / 180.0
	centerAltRad := centerAlt * math.Pi / 180.0
	centerAzRad := centerAz * math.Pi / 180.0
	fovRad := fov * math.Pi / 180.0

	// Convert to Cartesian coordinates on unit sphere
	starX := math.Cos(altRad) * math.Sin(azRad)
	starY := math.Cos(altRad) * math.Cos(azRad)
	starZ := math.Sin(altRad)

	// Center direction vector
	centerX := math.Cos(centerAltRad) * math.Sin(centerAzRad)
	centerY := math.Cos(centerAltRad) * math.Cos(centerAzRad)
	centerZ := math.Sin(centerAltRad)

	// Calculate angular separation (dot product)
	dotProduct := starX*centerX + starY*centerY + starZ*centerZ
	angularSep := math.Acos(math.Max(-1.0, math.Min(1.0, dotProduct)))

	// Check if within FOV
	if angularSep > fovRad/2.0 {
		return 0, 0, false
	}

	// Create coordinate system for projection
	// Up vector (celestial north)
	upX := 0.0
	upY := 0.0
	upZ := 1.0

	// Right vector (cross product of center and up)
	rightX := centerY*upZ - centerZ*upY
	rightY := centerZ*upX - centerX*upZ
	rightZ := centerX*upY - centerY*upX
	rightLen := math.Sqrt(rightX*rightX + rightY*rightY + rightZ*rightZ)
	if rightLen > 0 {
		rightX /= rightLen
		rightY /= rightLen
		rightZ /= rightLen
	}

	// Recalculate up vector (cross product of right and center)
	upX = rightY*centerZ - rightZ*centerY
	upY = rightZ*centerX - rightX*centerZ
	upZ = rightX*centerY - rightY*centerX

	// Project star onto tangent plane
	projRight := starX*rightX + starY*rightY + starZ*rightZ
	projUp := starX*upX + starY*upY + starZ*upZ

	// Scale by FOV
	scale := 2.0 / math.Tan(fovRad/2.0)
	screenX := projRight * scale
	screenY := projUp * scale

	// Convert to screen coordinates
	x = int(float64(screenWidth)/2.0 + screenX*float64(screenWidth)/2.0)
	y = int(float64(screenHeight)/2.0 - screenY*float64(screenHeight)/2.0)

	// Check bounds
	if x < 0 || x >= screenWidth || y < 0 || y >= screenHeight {
		return 0, 0, false
	}

	return x, y, true
}
