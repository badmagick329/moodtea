package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func bar(n int) string {
	// n is 1..5
	return strings.Repeat("█", n) + strings.Repeat("░", 5-n)
}

func renderLegend(title string, palette []string) string {
	const cellWidth = 4
	labelStyle := lipgloss.NewStyle().Bold(true)
	cellStyle := lipgloss.NewStyle().Width(cellWidth).Align(lipgloss.Center).Foreground(lipgloss.Color("#FFFFFF"))

	var b strings.Builder
	b.WriteString(labelStyle.Render(title))
	b.WriteString(" ")
	for i, c := range palette {
		b.WriteString(cellStyle.Background(lipgloss.Color(c)).Render(fmt.Sprintf("%d", i+1)))
	}
	return b.String()
}

func (m Model) View() tea.View {
	if m.err != nil {
		return tea.NewView("Error: " + m.err.Error() + "\n\n(q to quit)\n")
	}
	if len(m.days) == 0 {
		return tea.NewView("No data.\n\n(q to quit)\n")
	}

	var b strings.Builder
	start, daysInMonth, startWeekday := monthBounds(m.days)
	dayMap := buildDayMap(m.days)
	selected := m.days[m.cursor]

	b.WriteString("MoodTea — ")
	b.WriteString(start.Format("January 2006"))
	b.WriteString("\n")
	b.WriteString("←/→ (h/l), ↑/↓ (k/j) to move, q to quit\n\n")

	b.WriteString(renderCalendar("Mood", daysInMonth, startWeekday, dayMap, selected.Date.Day(), func(d dayInfo) int {
		return d.mood
	}, moodPalette))
	b.WriteString("\n")
	b.WriteString(renderLegend("Mood scale:", moodPalette))
	b.WriteString("\n\n")
	b.WriteString(renderCalendar("Energy", daysInMonth, startWeekday, dayMap, selected.Date.Day(), func(d dayInfo) int {
		return d.energy
	}, energyPalette))
	b.WriteString("\n")
	b.WriteString(renderLegend("Energy scale:", energyPalette))
	b.WriteString("\n")

	b.WriteString("Selected: ")
	b.WriteString(selected.Date.Format("2006-01-02"))
	b.WriteString("  mood ")
	b.WriteString(bar(selected.Mood))
	b.WriteString(fmt.Sprintf(" (%d)  energy ", selected.Mood))
	b.WriteString(bar(selected.Energy))
	b.WriteString(fmt.Sprintf(" (%d)\n", selected.Energy))

	return tea.NewView(b.String())
}
