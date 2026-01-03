# skyterm

A terminal-based astronomy application built as a Go TUI using the Bubbletea framework.

## Quick Start

```bash
# Build
go build -o skyterm ./cmd/skyterm

# Run
./skyterm
```

## Keybindings

### Navigation
- `↑/k` - Pan up
- `↓/j` - Pan down
- `←/h` - Pan left
- `→/l` - Pan right
- `K/J/H/L` - Fast pan

### Zoom
- `+` - Zoom in
- `-` - Zoom out
- `0` - Reset view

### Cardinal Directions
- `n` - North
- `s` - South
- `e` - East
- `w` - West
- `z` - Zenith

### Display Toggles
- `g` - Toggle coordinate grid
- `C` - Toggle constellation lines
- `N` - Toggle constellation names
- `p` - Toggle planets (Sun, Moon, planets)
- `P` - Toggle planet labels
- `d` - Toggle deep sky objects (Messier catalog)
- `S` - Toggle star labels (bright stars)
- `m` - Cycle magnitude limit

### Object Interaction
- `Enter` - Select nearest object to center
- `i` - Toggle info panel for selected object
- `c` - Center view on selected object
- `f` - Follow selected object (locks view)
- `/` - Search for object by name

### Time Controls
- `Space` - Pause/resume time flow
- `[` / `]` - Step time backward/forward
- `{` / `}` - Fast step time backward/forward (10x)
- `T` - Jump to current time (now)
- `t` - Set custom time

### General
- `?` - Help
- `q` - Quit

## Configuration

skyterm uses XDG-compliant configuration:
- Config file: `~/.config/skyterm/config.yaml` (or `$XDG_CONFIG_HOME/skyterm/config.yaml`)

Example configuration:
```yaml
location:
  latitude: 51.5074    # London
  longitude: -0.1278
  altitude: 11
  name: "London, UK"

display:
  magnitude_limit: 5.0
  show_constellation_lines: true
  show_constellation_names: false
  show_coordinate_grid: false
  show_planet_labels: true
  color_stars_by_type: true

time:
  use_utc: false
  time_step: "1m"    # Time step for [ and ] keys (e.g., "1m", "1h", "24h")
```

To change your city, edit the config file and update the `location` section with your coordinates.

**Note**: Default location is New York City (40.7°N, 74.0°W). The time controls allow you to pause, step through time, or jump to specific moments for observing celestial events.

## Development

```bash
# Run tests
go test ./...

# Run specific package tests
go test ./internal/astro

# Build for release
go build -ldflags="-s -w" -o skyterm ./cmd/skyterm
```

## License

TBD

## Credits

Created by Craig Derington
