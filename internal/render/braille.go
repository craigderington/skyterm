package render

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/craigderington/skyterm/internal/catalog"
)

// BrailleCanvas represents a high-resolution canvas using Braille Unicode characters
// Each Braille character is a 2x4 dot matrix, giving us 2x horizontal and 4x vertical resolution
type BrailleCanvas struct {
	Width     int // Width in Braille characters
	Height    int // Height in Braille characters
	PixelData [][]uint8
	Styles    [][]lipgloss.Style
}

// Braille Unicode base: U+2800
// Each dot in the 2x4 matrix corresponds to a bit:
//   1  8
//   2  16
//   4  32
//   64 128
const brailleBase = 0x2800

var brailleDots = [4][2]uint8{
	{1, 8},   // Row 0
	{2, 16},  // Row 1
	{4, 32},  // Row 2
	{64, 128}, // Row 3
}

// NewBrailleCanvas creates a new Braille canvas
func NewBrailleCanvas(width, height int) *BrailleCanvas {
	pixelData := make([][]uint8, height)
	styles := make([][]lipgloss.Style, height)
	for i := range pixelData {
		pixelData[i] = make([]uint8, width)
		styles[i] = make([]lipgloss.Style, width)
		for j := range styles[i] {
			styles[i][j] = lipgloss.NewStyle()
		}
	}

	return &BrailleCanvas{
		Width:     width,
		Height:    height,
		PixelData: pixelData,
		Styles:    styles,
	}
}

// SetPixel sets a pixel in the Braille canvas
// x, y are in pixel coordinates (2x width, 4x height of character grid)
func (bc *BrailleCanvas) SetPixel(x, y int, style lipgloss.Style) {
	if x < 0 || y < 0 {
		return
	}

	charX := x / 2
	charY := y / 4

	if charX >= bc.Width || charY >= bc.Height {
		return
	}

	dotX := x % 2
	dotY := y % 4

	// Set the appropriate bit
	bc.PixelData[charY][charX] |= brailleDots[dotY][dotX]
	bc.Styles[charY][charX] = style
}

// Render converts the Braille canvas to a string
func (bc *BrailleCanvas) Render() string {
	var result string
	for y := 0; y < bc.Height; y++ {
		for x := 0; x < bc.Width; x++ {
			brailleChar := rune(brailleBase + int(bc.PixelData[y][x]))
			result += bc.Styles[y][x].Render(string(brailleChar))
		}
		if y < bc.Height-1 {
			result += "\n"
		}
	}
	return result
}

// Clear clears the Braille canvas
func (bc *BrailleCanvas) Clear() {
	for y := 0; y < bc.Height; y++ {
		for x := 0; x < bc.Width; x++ {
			bc.PixelData[y][x] = 0
			bc.Styles[y][x] = lipgloss.NewStyle()
		}
	}
}

// RenderStarsBraille renders stars using Braille characters for higher resolution
func RenderStarsBraille(bc *BrailleCanvas, stars []catalog.Star, centerAlt, centerAz, fov, magLimit float64) {
	// Pixel dimensions (2x width, 4x height)
	pixelWidth := bc.Width * 2
	pixelHeight := bc.Height * 4

	for _, star := range stars {
		if star.Magnitude > magLimit {
			continue
		}

		// Project to normalized coordinates
		x, y, visible := projectNormalized(star.Altitude, star.Azimuth, centerAlt, centerAz, fov)
		if !visible {
			continue
		}

		// Convert to pixel coordinates
		px := int((x + 1.0) * float64(pixelWidth) / 2.0)
		py := int((1.0 - y) * float64(pixelHeight) / 2.0)

		// Get color for spectral type
		color := getColorForSpectralType(star.SpectralType)
		style := lipgloss.NewStyle().Foreground(color)

		// For brighter stars, set multiple pixels
		if star.Magnitude < 1.0 {
			// Very bright - 3x3 pixel cluster
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					bc.SetPixel(px+dx, py+dy, style)
				}
			}
		} else if star.Magnitude < 3.0 {
			// Bright - 2x2 pixel cluster
			for dy := 0; dy <= 1; dy++ {
				for dx := 0; dx <= 1; dx++ {
					bc.SetPixel(px+dx, py+dy, style)
				}
			}
		} else {
			// Normal - single pixel
			bc.SetPixel(px, py, style)
		}
	}
}

// projectNormalized returns normalized coordinates in range [-1, 1]
func projectNormalized(alt, az, centerAlt, centerAz, fov float64) (x, y float64, visible bool) {
	// Use existing project function but normalize output
	// This is a simplified version - we'd need to extract the math from project()
	// For now, return placeholder
	return 0, 0, false
}
