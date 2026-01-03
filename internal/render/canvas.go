package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Cell struct {
	Char  rune
	Style lipgloss.Style
}

type Canvas struct {
	Width  int
	Height int
	Cells  [][]Cell
}

func NewCanvas(width, height int) *Canvas {
	cells := make([][]Cell, height)
	for i := range cells {
		cells[i] = make([]Cell, width)
		for j := range cells[i] {
			cells[i][j] = Cell{
				Char:  ' ',
				Style: lipgloss.NewStyle(),
			}
		}
	}

	return &Canvas{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

func (c *Canvas) Clear() {
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			c.Cells[y][x] = Cell{
				Char:  ' ',
				Style: lipgloss.NewStyle(),
			}
		}
	}
}

func (c *Canvas) Set(x, y int, char rune, style lipgloss.Style) {
	if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
		c.Cells[y][x] = Cell{
			Char:  char,
			Style: style,
		}
	}
}

func (c *Canvas) Render() string {
	var sb strings.Builder
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			cell := c.Cells[y][x]
			sb.WriteString(cell.Style.Render(string(cell.Char)))
		}
		if y < c.Height-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}
