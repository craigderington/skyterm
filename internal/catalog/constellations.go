package catalog

// ConstellationLine represents a line connecting stars in a constellation
type ConstellationLine struct {
	Star1Name string
	Star2Name string
}

// Constellation represents a constellation with its lines and metadata
type Constellation struct {
	Name       string
	Abbreviation string
	Lines      []ConstellationLine
}

// GetConstellations returns all constellation definitions
func GetConstellations() []Constellation {
	return []Constellation{
		{
			Name:         "Orion",
			Abbreviation: "Ori",
			Lines: []ConstellationLine{
				// Shoulders to belt
				{"Betelgeuse", "Bellatrix"},
				{"Betelgeuse", "Alnitak"},
				{"Bellatrix", "Mintaka"},

				// Belt
				{"Alnitak", "Alnilam"},
				{"Alnilam", "Mintaka"},

				// Belt to legs
				{"Alnitak", "Saiph"},
				{"Mintaka", "Rigel"},

				// Legs
				{"Rigel", "Saiph"},
			},
		},
		{
			Name:         "Ursa Major",
			Abbreviation: "UMa",
			Lines: []ConstellationLine{
				// Big Dipper bowl
				{"Dubhe", "Merak"},
				{"Merak", "Phecda"},
				{"Phecda", "Megrez"},
				{"Megrez", "Dubhe"},

				// Big Dipper handle
				{"Megrez", "Alioth"},
				{"Alioth", "Mizar"},
				{"Mizar", "Alkaid"},
			},
		},
		{
			Name:         "Summer Triangle",
			Abbreviation: "---",
			Lines: []ConstellationLine{
				{"Vega", "Deneb"},
				{"Deneb", "Altair"},
				{"Altair", "Vega"},
			},
		},
		{
			Name:         "Crux",
			Abbreviation: "Cru",
			Lines: []ConstellationLine{
				// Southern Cross
				{"Acrux", "Gacrux"},
				{"Mimosa", "Gacrux"},
			},
		},
		{
			Name:         "Canis Major",
			Abbreviation: "CMa",
			Lines: []ConstellationLine{
				// Simple representation with Sirius
				{"Sirius", "Sirius"}, // Placeholder - would need more stars
			},
		},
		{
			Name:         "Leo",
			Abbreviation: "Leo",
			Lines: []ConstellationLine{
				// Simple representation with Regulus
				{"Regulus", "Regulus"}, // Placeholder
			},
		},
		{
			Name:         "Gemini",
			Abbreviation: "Gem",
			Lines: []ConstellationLine{
				{"Castor", "Pollux"},
			},
		},
	}
}

// ConstellationLabel represents where to place a constellation label
type ConstellationLabel struct {
	Name     string
	StarName string // Label near this star
}

// GetConstellationLabels returns label positions for constellations
func GetConstellationLabels() []ConstellationLabel {
	return []ConstellationLabel{
		{Name: "ORION", StarName: "Alnilam"},
		{Name: "URSA MAJOR", StarName: "Alioth"},
		{Name: "CRUX", StarName: "Acrux"},
		{Name: "GEMINI", StarName: "Castor"},
		{Name: "LEO", StarName: "Regulus"},
		{Name: "CANIS MAJOR", StarName: "Sirius"},
	}
}
