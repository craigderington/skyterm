package app

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	// Navigation
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	FastUp    key.Binding
	FastDown  key.Binding
	FastLeft  key.Binding
	FastRight key.Binding

	// Zoom
	ZoomIn  key.Binding
	ZoomOut key.Binding
	Reset   key.Binding

	// Cardinal directions
	North  key.Binding
	South  key.Binding
	East   key.Binding
	West   key.Binding
	Zenith key.Binding

	// Display toggles
	Grid           key.Binding
	Constellations key.Binding
	Names          key.Binding
	Planets        key.Binding
	PlanetLabels   key.Binding
	DeepSky        key.Binding
	StarLabels     key.Binding
	Magnitude      key.Binding

	// Selection and interaction
	Select     key.Binding
	Info       key.Binding
	ViewImage  key.Binding
	Center     key.Binding
	Follow     key.Binding
	Search     key.Binding

	// Time controls
	PauseResume    key.Binding
	StepBackward   key.Binding
	StepForward    key.Binding
	FastStepBack   key.Binding
	FastStepForward key.Binding
	JumpToNow      key.Binding
	SetTime        key.Binding

	// General
	Help      key.Binding
	Quit      key.Binding
	CloseHelp key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "pan up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "pan down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "pan left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "pan right"),
		),
		FastUp: key.NewBinding(
			key.WithKeys("K"),
			key.WithHelp("K", "fast pan up"),
		),
		FastDown: key.NewBinding(
			key.WithKeys("J"),
			key.WithHelp("J", "fast pan down"),
		),
		FastLeft: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "fast pan left"),
		),
		FastRight: key.NewBinding(
			key.WithKeys("L"),
			key.WithHelp("L", "fast pan right"),
		),

		// Zoom
		ZoomIn: key.NewBinding(
			key.WithKeys("+", "="),
			key.WithHelp("+", "zoom in"),
		),
		ZoomOut: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "zoom out"),
		),
		Reset: key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "reset view"),
		),

		// Cardinals
		North: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "north"),
		),
		South: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "south"),
		),
		East: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "east"),
		),
		West: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "west"),
		),
		Zenith: key.NewBinding(
			key.WithKeys("z"),
			key.WithHelp("z", "zenith"),
		),

		// Display
		Grid: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "toggle grid"),
		),
		Constellations: key.NewBinding(
			key.WithKeys("C"),
			key.WithHelp("C", "toggle constellations"),
		),
		Names: key.NewBinding(
			key.WithKeys("N"),
			key.WithHelp("N", "toggle names"),
		),
		Planets: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "toggle planets"),
		),
		PlanetLabels: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle planet labels"),
		),
		DeepSky: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "toggle deep sky objects"),
		),
		StarLabels: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle star labels"),
		),
		Magnitude: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "cycle magnitude"),
		),

		// Selection and interaction
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select object"),
		),
		Info: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "toggle info"),
		),
		ViewImage: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "view image"),
		),
		Center: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "center on selected"),
		),
		Follow: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "follow selected"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search objects"),
		),

		// Time controls
		PauseResume: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "pause/resume time"),
		),
		StepBackward: key.NewBinding(
			key.WithKeys("["),
			key.WithHelp("[", "step time backward"),
		),
		StepForward: key.NewBinding(
			key.WithKeys("]"),
			key.WithHelp("]", "step time forward"),
		),
		FastStepBack: key.NewBinding(
			key.WithKeys("{"),
			key.WithHelp("{", "fast step backward"),
		),
		FastStepForward: key.NewBinding(
			key.WithKeys("}"),
			key.WithHelp("}", "fast step forward"),
		),
		JumpToNow: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "jump to current time"),
		),
		SetTime: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "set time"),
		),

		// General
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		CloseHelp: key.NewBinding(
			key.WithKeys("esc", "?"),
			key.WithHelp("esc/?", "close help"),
		),
	}
}
