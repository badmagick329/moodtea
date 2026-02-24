package ui

import "charm.land/lipgloss/v2"

var (
	moodPalette = []string{
		"#0B2A14",
		"#1E5A2C",
		"#2E7A3F",
		"#3F9A52",
		"#5FB86A",
	}
	energyPalette = []string{
		"#2B1D00",
		"#6B3D00",
		"#B65A00",
		"#E37B00",
		"#FFB000",
	}

	headerCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4A5568")).
			Padding(0, 1)
	footerCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4A5568")).
			Padding(0, 1)
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F7FAFC"))
	subtleTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))
	warnTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F6AD55"))
	chipStyle       = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("#303642")).
			Foreground(lipgloss.Color("#CBD5E1"))
	barFilledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F8FAFC"))
	barEmptyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	dividerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#2D3748"))
)
