package app

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/craigderington/skyterm/internal/astro"
	"github.com/craigderington/skyterm/internal/catalog"
	"github.com/craigderington/skyterm/internal/config"
	"github.com/craigderington/skyterm/internal/render"
	"github.com/craigderington/skyterm/internal/ui"
)

type Model struct {
	keys   keyMap
	width  int
	height int

	// View parameters
	altitude float64 // Altitude (elevation) in degrees
	azimuth  float64 // Azimuth in degrees
	fov      float64 // Field of view in degrees

	// Display options
	showGrid           bool
	showConstellations bool
	showNames          bool
	showPlanets        bool
	showPlanetLabels   bool
	showDeepSky        bool
	showStarLabels     bool
	magnitudeLimit     float64
	showHelp           bool
	showInfo           bool

	// Interaction state
	selectedObject *SelectedObject
	following      bool
	searchMode     bool
	searchQuery    string
	timeInputMode  bool
	timeInput      string

	// Time and location
	currentTime    time.Time
	observer       *astro.Observer
	paused         bool
	timeStep       time.Duration
	timeMultiplier float64
	realTimeBase   time.Time // When we started or last resumed real-time mode

	// Data
	starCatalog     *catalog.StarCatalog
	deepSkyCatalog  *catalog.DeepSkyCatalog
	planetarySystem *astro.PlanetarySystem
	canvas          *render.Canvas

	// Config
	config *config.Config
}

func New() Model {
	cfg, _ := config.Load()

	// Parse time step from config
	timeStep, err := time.ParseDuration(cfg.Time.TimeStep)
	if err != nil {
		timeStep = 1 * time.Minute // Default to 1 minute
	}

	now := time.Now()

	return Model{
		keys:               newKeyMap(),
		altitude:           45.0,  // Start looking 45° up
		azimuth:            180.0, // South
		fov:                60.0,  // 60° field of view
		showGrid:           cfg.Display.ShowCoordinateGrid,
		showConstellations: cfg.Display.ShowConstellationLines,
		showNames:          cfg.Display.ShowConstellationNames,
		showPlanets:        false,
		showPlanetLabels:   cfg.Display.ShowPlanetLabels,
		showDeepSky:        false,
		showStarLabels:     true, // Show star labels by default
		magnitudeLimit:     cfg.Display.MagnitudeLimit,
		currentTime:        now,
		observer:           cfg.Observer(),
		paused:             false,
		timeStep:           timeStep,
		timeMultiplier:     1.0,
		realTimeBase:       now,
		starCatalog:        catalog.NewStarCatalog(),
		deepSkyCatalog:     catalog.NewDeepSkyCatalog(),
		planetarySystem:    &astro.PlanetarySystem{},
		config:             cfg,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		realNow := time.Time(msg)

		// Update current time based on pause state
		if !m.paused {
			// Calculate elapsed time since base with multiplier applied
			elapsed := realNow.Sub(m.realTimeBase)
			scaledElapsed := time.Duration(float64(elapsed) * m.timeMultiplier)
			m.currentTime = m.realTimeBase.Add(scaledElapsed)

			// Update base for next tick if running at normal speed
			if m.timeMultiplier == 1.0 {
				m.realTimeBase = realNow
			}
		}
		// If paused, currentTime stays frozen

		m.starCatalog.UpdatePositions(m.observer, m.currentTime)
		m.deepSkyCatalog.UpdatePositions(m.observer, m.currentTime)
		m.planetarySystem = astro.CalculatePlanets(m.currentTime, m.observer)

		// Update following if active
		m.UpdateFollowing()

		return m, tickCmd()

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 2 // Reserve space for status bar
		m.canvas = render.NewCanvas(m.width, m.height)
		return m, nil

	case tea.KeyMsg:
		// Handle time input mode
		if m.timeInputMode {
			switch msg.String() {
			case "esc":
				m.timeInputMode = false
				m.timeInput = ""
				return m, nil
			case "enter":
				// Parse and set time
				if m.timeInput != "" {
					if parsedTime, err := time.Parse("2006-01-02 15:04:05", m.timeInput); err == nil {
						m.currentTime = parsedTime
						m.realTimeBase = time.Now()
						m.paused = true // Pause when manually setting time
					} else if parsedTime, err := time.Parse("2006-01-02 15:04", m.timeInput); err == nil {
						m.currentTime = parsedTime
						m.realTimeBase = time.Now()
						m.paused = true
					} else if parsedTime, err := time.Parse("2006-01-02", m.timeInput); err == nil {
						m.currentTime = parsedTime
						m.realTimeBase = time.Now()
						m.paused = true
					}
				}
				m.timeInputMode = false
				m.timeInput = ""
				return m, nil
			case "backspace":
				if len(m.timeInput) > 0 {
					m.timeInput = m.timeInput[:len(m.timeInput)-1]
				}
				return m, nil
			default:
				// Add character to input
				if len(msg.String()) == 1 {
					m.timeInput += msg.String()
				}
				return m, nil
			}
		}

		// Handle search mode
		if m.searchMode {
			switch msg.String() {
			case "esc":
				m.searchMode = false
				m.searchQuery = ""
				return m, nil
			case "enter":
				m.performSearch()
				m.searchMode = false
				m.searchQuery = ""
				return m, nil
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
				return m, nil
			default:
				// Add character to query
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
				}
				return m, nil
			}
		}

		// Handle help screen
		if m.showHelp {
			if key.Matches(msg, m.keys.CloseHelp) {
				m.showHelp = false
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.Help):
			m.showHelp = true
			return m, nil

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		// Selection and interaction
		case key.Matches(msg, m.keys.Select):
			m.SelectNearestObject()
			return m, nil

		case key.Matches(msg, m.keys.Info):
			m.showInfo = !m.showInfo
			return m, nil

		case key.Matches(msg, m.keys.Center):
			m.CenterOnSelected()
			return m, nil

		case key.Matches(msg, m.keys.Follow):
			if m.selectedObject != nil {
				m.following = !m.following
			}
			return m, nil

		case key.Matches(msg, m.keys.Search):
			m.searchMode = true
			m.searchQuery = ""
			return m, nil

		// Time controls
		case key.Matches(msg, m.keys.PauseResume):
			m.paused = !m.paused
			if !m.paused {
				// Resuming - reset real time base
				m.realTimeBase = time.Now()
			}
			return m, nil

		case key.Matches(msg, m.keys.StepForward):
			m.currentTime = m.currentTime.Add(m.timeStep)
			m.paused = true
			return m, nil

		case key.Matches(msg, m.keys.StepBackward):
			m.currentTime = m.currentTime.Add(-m.timeStep)
			m.paused = true
			return m, nil

		case key.Matches(msg, m.keys.FastStepForward):
			m.currentTime = m.currentTime.Add(m.timeStep * 10)
			m.paused = true
			return m, nil

		case key.Matches(msg, m.keys.FastStepBack):
			m.currentTime = m.currentTime.Add(-m.timeStep * 10)
			m.paused = true
			return m, nil

		case key.Matches(msg, m.keys.JumpToNow):
			m.currentTime = time.Now()
			m.realTimeBase = time.Now()
			m.paused = false
			return m, nil

		case key.Matches(msg, m.keys.SetTime):
			m.timeInputMode = true
			m.timeInput = ""
			return m, nil

		// Navigation
		case key.Matches(msg, m.keys.Up):
			m.altitude = math.Min(90.0, m.altitude+m.config.Controls.PanSpeed)
		case key.Matches(msg, m.keys.Down):
			m.altitude = math.Max(-90.0, m.altitude-m.config.Controls.PanSpeed)
		case key.Matches(msg, m.keys.Left):
			m.azimuth = math.Mod(m.azimuth-m.config.Controls.PanSpeed+360, 360)
		case key.Matches(msg, m.keys.Right):
			m.azimuth = math.Mod(m.azimuth+m.config.Controls.PanSpeed, 360)

		// Fast navigation
		case key.Matches(msg, m.keys.FastUp):
			m.altitude = math.Min(90.0, m.altitude+m.config.Controls.PanSpeed*m.config.Controls.FastPanMultiplier)
		case key.Matches(msg, m.keys.FastDown):
			m.altitude = math.Max(-90.0, m.altitude-m.config.Controls.PanSpeed*m.config.Controls.FastPanMultiplier)
		case key.Matches(msg, m.keys.FastLeft):
			m.azimuth = math.Mod(m.azimuth-m.config.Controls.PanSpeed*m.config.Controls.FastPanMultiplier+360, 360)
		case key.Matches(msg, m.keys.FastRight):
			m.azimuth = math.Mod(m.azimuth+m.config.Controls.PanSpeed*m.config.Controls.FastPanMultiplier, 360)

		// Zoom
		case key.Matches(msg, m.keys.ZoomIn):
			m.fov = math.Max(10.0, m.fov/m.config.Controls.ZoomStep)
		case key.Matches(msg, m.keys.ZoomOut):
			m.fov = math.Min(120.0, m.fov*m.config.Controls.ZoomStep)

		// Reset
		case key.Matches(msg, m.keys.Reset):
			m.altitude = 45.0
			m.azimuth = 180.0
			m.fov = 60.0

		// Cardinals
		case key.Matches(msg, m.keys.North):
			m.azimuth = 0.0
		case key.Matches(msg, m.keys.South):
			m.azimuth = 180.0
		case key.Matches(msg, m.keys.East):
			m.azimuth = 90.0
		case key.Matches(msg, m.keys.West):
			m.azimuth = 270.0
		case key.Matches(msg, m.keys.Zenith):
			m.altitude = 90.0

		// Display toggles
		case key.Matches(msg, m.keys.Grid):
			m.showGrid = !m.showGrid
		case key.Matches(msg, m.keys.Constellations):
			m.showConstellations = !m.showConstellations
		case key.Matches(msg, m.keys.Names):
			m.showNames = !m.showNames
		case key.Matches(msg, m.keys.Planets):
			m.showPlanets = !m.showPlanets
		case key.Matches(msg, m.keys.PlanetLabels):
			m.showPlanetLabels = !m.showPlanetLabels
		case key.Matches(msg, m.keys.DeepSky):
			m.showDeepSky = !m.showDeepSky
		case key.Matches(msg, m.keys.StarLabels):
			m.showStarLabels = !m.showStarLabels
		case key.Matches(msg, m.keys.Magnitude):
			// Cycle through magnitude limits: 3, 4, 5, 6
			switch m.magnitudeLimit {
			case 3.0:
				m.magnitudeLimit = 4.0
			case 4.0:
				m.magnitudeLimit = 5.0
			case 5.0:
				m.magnitudeLimit = 6.0
			default:
				m.magnitudeLimit = 3.0
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// Show time input modal if in time input mode
	if m.timeInputMode {
		return ui.RenderTimeInput(m.timeInput, m.width, m.height+2)
	}

	// Show search box if in search mode
	if m.searchMode {
		return ui.RenderSearchBox(m.searchQuery, m.width, m.height+2)
	}

	// Show help screen if requested
	if m.showHelp {
		return ui.RenderHelp(m.width, m.height+2)
	}

	// Clear canvas
	m.canvas.Clear()

	// Render coordinate grid (if enabled)
	if m.showGrid {
		render.RenderGrid(m.canvas, m.altitude, m.azimuth, m.fov)
	}

	// Render constellation lines (if enabled)
	if m.showConstellations {
		render.RenderConstellations(
			m.canvas,
			m.starCatalog.Stars(),
			catalog.GetConstellations(),
			m.altitude,
			m.azimuth,
			m.fov,
		)
	}

	// Render deep sky objects (if enabled)
	if m.showDeepSky {
		render.RenderDeepSkyObjects(
			m.canvas,
			m.deepSkyCatalog.Objects(),
			m.altitude,
			m.azimuth,
			m.fov,
			m.magnitudeLimit,
		)
		render.RenderDeepSkyLabels(
			m.canvas,
			m.deepSkyCatalog.Objects(),
			m.altitude,
			m.azimuth,
			m.fov,
			m.magnitudeLimit,
		)
	}

	// Render stars
	render.RenderStars(m.canvas, m.starCatalog.Stars(), m.altitude, m.azimuth, m.fov, m.magnitudeLimit)

	// Render star labels (if enabled)
	if m.showStarLabels {
		render.RenderStarLabels(m.canvas, m.starCatalog.Stars(), m.altitude, m.azimuth, m.fov, m.magnitudeLimit)
	}

	// Render planets (if enabled)
	if m.showPlanets {
		render.RenderPlanets(m.canvas, m.planetarySystem, m.altitude, m.azimuth, m.fov)
	}

	// Render planet labels (if enabled)
	if m.showPlanetLabels {
		render.RenderPlanetLabels(m.canvas, m.planetarySystem, m.altitude, m.azimuth, m.fov)
	}

	// Render constellation names (if enabled)
	if m.showNames {
		render.RenderConstellationLabels(
			m.canvas,
			m.starCatalog.Stars(),
			catalog.GetConstellationLabels(),
			m.altitude,
			m.azimuth,
			m.fov,
		)
	}

	// Build the view
	skyView := m.canvas.Render()
	statusBar := m.renderStatusBar()

	view := lipgloss.JoinVertical(lipgloss.Left, skyView, statusBar)

	// Overlay info panel if requested
	if m.showInfo && m.selectedObject != nil {
		// Convert to ui.ObjectInfo
		info := &ui.ObjectInfo{
			Type:    m.selectedObject.Type,
			Name:    m.selectedObject.Name,
			Star:    m.selectedObject.Star,
			Planet:  m.selectedObject.Planet,
			DeepSky: m.selectedObject.DeepSky,
		}
		infoPanel := ui.RenderInfoPanel(info, m.observer, m.width, m.height+2)
		// Overlay the panel on top of the view
		view = lipgloss.Place(
			m.width,
			m.height+2,
			lipgloss.Left,
			lipgloss.Top,
			view,
		)
		view = view + "\n" + infoPanel
	}

	return view
}

func (m Model) renderStatusBar() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("235"))

	// Format current time
	timeStr := m.currentTime.Format("2006-01-02 15:04:05 MST")
	if m.config.Time.UseUTC {
		timeStr = m.currentTime.UTC().Format("2006-01-02 15:04:05 UTC")
	}

	// Add paused indicator
	pausedIndicator := ""
	if m.paused {
		pausedIndicator = " [PAUSED]"
	}

	// Get sidereal time
	lst := m.observer.LST(m.currentTime)
	lstStr := astro.FormatSiderealTime(lst)

	// Build toggle indicators
	toggles := ""
	if m.showGrid {
		toggles += "G"
	}
	if m.showConstellations {
		toggles += "C"
	}
	if m.showNames {
		toggles += "N"
	}
	if m.showPlanets {
		toggles += "p"
	}
	if m.showPlanetLabels {
		toggles += "P"
	}
	if m.showDeepSky {
		toggles += "D"
	}
	if m.showStarLabels {
		toggles += "S"
	}
	if toggles != "" {
		toggles = " [" + toggles + "]"
	}

	// Build status bar sections
	left := fmt.Sprintf(" Alt: %.1f° Az: %.1f° │ FOV: %.1f°%s", m.altitude, m.azimuth, m.fov, toggles)
	center := fmt.Sprintf(" %s%s │ LST: %s", timeStr, pausedIndicator, lstStr)
	right := fmt.Sprintf("Mag: %.1f │ %s ", m.magnitudeLimit, m.observer.Name)

	// Calculate padding
	usedWidth := lipgloss.Width(left) + lipgloss.Width(center) + lipgloss.Width(right)
	padding := m.width - usedWidth
	if padding < 0 {
		// If too wide, simplify
		left = fmt.Sprintf(" %.1f°/%.1f°%s", m.altitude, m.azimuth, toggles)
		right = fmt.Sprintf("%.1f ", m.magnitudeLimit)
		usedWidth = lipgloss.Width(left) + lipgloss.Width(center) + lipgloss.Width(right)
		padding = m.width - usedWidth
	}
	if padding < 0 {
		padding = 0
	}

	leftBar := style.Render(left)
	centerBar := style.Render(center)
	rightBar := style.Render(right)

	// Split padding between two gaps
	leftPad := padding / 2
	rightPad := padding - leftPad

	leftSpacer := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Render(lipgloss.PlaceHorizontal(leftPad, lipgloss.Left, ""))

	rightSpacer := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Render(lipgloss.PlaceHorizontal(rightPad, lipgloss.Left, ""))

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBar, leftSpacer, centerBar, rightSpacer, rightBar)
}
