package image

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

// TerminalCapability represents what image rendering the terminal supports
type TerminalCapability int

const (
	CapabilityNone TerminalCapability = iota
	CapabilityUnicode
	CapabilitySixel
	CapabilityKitty
	CapabilityITerm2
)

// TerminalRenderer handles rendering images in the terminal
type TerminalRenderer struct {
	capability TerminalCapability
}

// NewTerminalRenderer creates a new terminal image renderer
func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{
		capability: detectTerminalCapability(),
	}
}

// detectTerminalCapability detects what image protocols the terminal supports
func detectTerminalCapability() TerminalCapability {
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")
	kittyWindowID := os.Getenv("KITTY_WINDOW_ID")

	// Check for Kitty - it sets KITTY_WINDOW_ID environment variable
	if kittyWindowID != "" || strings.Contains(term, "kitty") || termProgram == "kitty" {
		return CapabilityKitty
	}

	// Check for iTerm2
	if termProgram == "iTerm.app" {
		return CapabilityITerm2
	}

	// Check for Sixel support (xterm, mlterm, etc.)
	if strings.Contains(term, "xterm") {
		// Could query terminal with escape sequences, but for now assume no sixel
		// return CapabilitySixel
	}

	// Fallback to Unicode block art
	return CapabilityUnicode
}

// Render renders an image for terminal display
func (r *TerminalRenderer) Render(imgData []byte, width, height int) (string, error) {
	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	switch r.capability {
	case CapabilityKitty:
		return r.renderKitty(imgData, width, height)
	case CapabilityITerm2:
		return r.renderITerm2(imgData, width, height)
	case CapabilitySixel:
		return r.renderSixel(img, width, height)
	default:
		return r.renderUnicode(img, width, height)
	}
}

// renderKitty renders image using Kitty graphics protocol
func (r *TerminalRenderer) renderKitty(imgData []byte, width, height int) (string, error) {
	// Kitty graphics protocol: https://sw.kovidgoyal.net/kitty/graphics-protocol/

	// Detect image format
	format := detectImageFormat(imgData)
	var formatCode int
	switch format {
	case "png":
		formatCode = 100
	case "jpeg", "jpg":
		formatCode = 24
	default:
		// Unknown format, try PNG
		formatCode = 100
	}

	// Encode image data
	encoded := base64.StdEncoding.EncodeToString(imgData)

	// Split into chunks of 4096 bytes (Kitty limitation)
	const chunkSize = 4096
	var output strings.Builder

	for i := 0; i < len(encoded); i += chunkSize {
		end := i + chunkSize
		if end > len(encoded) {
			end = len(encoded)
		}
		chunk := encoded[i:end]

		more := 0
		if end < len(encoded) {
			more = 1
		}

		if i == 0 {
			// First chunk with parameters
			// f=format, a=T (transmit and display), t=d (direct), m=more chunks
			// c=columns, r=rows (in cells)
			output.WriteString(fmt.Sprintf("\033_Gf=%d,a=T,t=d,c=%d,r=%d,m=%d;%s\033\\",
				formatCode, width, height, more, chunk))
		} else {
			// Continuation chunks
			output.WriteString(fmt.Sprintf("\033_Gm=%d;%s\033\\", more, chunk))
		}
	}

	// Add a newline after the image
	output.WriteString("\n")

	return output.String(), nil
}

// detectImageFormat detects the format of image data from magic bytes
func detectImageFormat(data []byte) string {
	if len(data) < 4 {
		return "unknown"
	}

	// PNG magic bytes: 89 50 4E 47
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "png"
	}

	// JPEG magic bytes: FF D8 FF
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "jpeg"
	}

	// GIF magic bytes: 47 49 46
	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
		return "gif"
	}

	return "unknown"
}

// renderITerm2 renders image using iTerm2 inline images protocol
func (r *TerminalRenderer) renderITerm2(imgData []byte, width, height int) (string, error) {
	// iTerm2 inline images: https://iterm2.com/documentation-images.html
	encoded := base64.StdEncoding.EncodeToString(imgData)

	return fmt.Sprintf("\033]1337;File=inline=1;width=%d;height=%d:%s\007",
		width, height, encoded), nil
}

// renderSixel renders image using Sixel protocol (TODO: implement)
func (r *TerminalRenderer) renderSixel(img image.Image, width, height int) (string, error) {
	// Sixel is complex, for now fallback to Unicode
	return r.renderUnicode(img, width, height)
}

// renderUnicode renders image as Unicode block art
func (r *TerminalRenderer) renderUnicode(img image.Image, width, height int) (string, error) {
	// Resize image to fit terminal cells
	// Each cell is roughly 2:1 ratio (height:width)
	resized := resize.Resize(uint(width), uint(height*2), img, resize.Lanczos3)

	var output strings.Builder

	// Use half-block characters for 2x vertical resolution
	for y := 0; y < height*2; y += 2 {
		for x := 0; x < width; x++ {
			// Get colors of two vertical pixels
			topColor := resized.At(x, y)
			bottomColor := resized.At(x, min(y+1, height*2-1))

			// Convert to block character
			char, fg, bg := blockChar(topColor, bottomColor)

			// Output with ANSI colors
			output.WriteString(ansiColor(char, fg, bg))
		}
		output.WriteString("\033[0m\n") // Reset and newline
	}

	return output.String(), nil
}

// blockChar determines the best Unicode block character and colors
func blockChar(topColor, bottomColor color.Color) (rune, color.Color, color.Color) {
	tr, tg, tb, ta := topColor.RGBA()
	br, bg, bb, ba := bottomColor.RGBA()

	// Convert to 8-bit
	tr, tg, tb, ta = tr>>8, tg>>8, tb>>8, ta>>8
	br, bg, bb, ba = br>>8, bg>>8, bb>>8, ba>>8

	// Calculate brightness
	topBright := (tr + tg + tb) / 3
	bottomBright := (br + bg + bb) / 3

	// If both transparent, use space
	if ta < 128 && ba < 128 {
		return ' ', topColor, bottomColor
	}

	// If top is transparent, use lower half block
	if ta < 128 {
		return '▄', bottomColor, color.Black
	}

	// If bottom is transparent, use upper half block
	if ba < 128 {
		return '▀', topColor, color.Black
	}

	// Similar colors, use full block or space
	if abs(int(topBright)-int(bottomBright)) < 30 {
		if topBright > 128 {
			return '█', topColor, topColor
		}
		return ' ', color.Black, color.Black
	}

	// Different colors, use half block
	if topBright > bottomBright {
		return '▀', topColor, bottomColor
	}
	return '▄', bottomColor, topColor
}

// ansiColor returns ANSI escape sequence for colored character
func ansiColor(char rune, fg, bg color.Color) string {
	fr, fg1, fb, _ := fg.RGBA()
	br, bg1, bb, _ := bg.RGBA()

	// Convert to 8-bit
	fr, fg1, fb = fr>>8, fg1>>8, fb>>8
	br, bg1, bb = br>>8, bg1>>8, bb>>8

	// Use 24-bit true color (RGB)
	return fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm%c",
		fr, fg1, fb, br, bg1, bb, char)
}

// abs returns absolute value of int
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// min returns minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RenderImageFromURL fetches and renders an image from a URL
func (r *TerminalRenderer) RenderImageFromURL(imageURL string, width, height int) (string, error) {
	// Use Wikipedia client to download
	client, err := NewWikipediaClient("")
	if err != nil {
		return "", err
	}

	imgData, err := client.DownloadImage(imageURL)
	if err != nil {
		return "", err
	}

	return r.Render(imgData, width, height)
}

// GetCapability returns the detected terminal capability
func (r *TerminalRenderer) GetCapability() TerminalCapability {
	return r.capability
}

// CapabilityString returns a human-readable string for the capability
func (r *TerminalRenderer) CapabilityString() string {
	switch r.capability {
	case CapabilityKitty:
		return "Kitty Graphics Protocol"
	case CapabilityITerm2:
		return "iTerm2 Inline Images"
	case CapabilitySixel:
		return "Sixel"
	case CapabilityUnicode:
		return "Unicode Block Art"
	default:
		return "None"
	}
}

// CalculateAspectRatio calculates appropriate dimensions maintaining aspect ratio
func CalculateAspectRatio(imgWidth, imgHeight, maxWidth, maxHeight int) (int, int) {
	aspectRatio := float64(imgWidth) / float64(imgHeight)

	// Terminal cells are roughly 2:1 (height:width), so adjust
	cellAspect := 2.0
	adjustedAspect := aspectRatio * cellAspect

	width := maxWidth
	height := int(math.Round(float64(width) / adjustedAspect))

	if height > maxHeight {
		height = maxHeight
		width = int(math.Round(float64(height) * adjustedAspect))
	}

	// Ensure minimum size
	if width < 10 {
		width = 10
	}
	if height < 5 {
		height = 5
	}

	return width, height
}
