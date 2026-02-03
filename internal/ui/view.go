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
	vm := m.buildViewModel()
	if !vm.HasData {
		return tea.NewView("No data.\n\n(q to quit)\n")
	}

	var b strings.Builder

	b.WriteString("MoodTea — ")
	b.WriteString(vm.MonthLabel)
	b.WriteString("\n")
	b.WriteString(helpLine)
	b.WriteString("\n\n")

	b.WriteString(renderCalendar("Mood", vm.DaysInMonth, vm.StartWeekday, vm.DayMap, vm.Selected.Date.Day(), func(d dayInfo) int {
		return d.mood
	}, moodPalette))
	b.WriteString("\n")
	b.WriteString(renderLegend("Mood scale:", moodPalette))
	b.WriteString("\n\n")
	b.WriteString(renderCalendar("Energy", vm.DaysInMonth, vm.StartWeekday, vm.DayMap, vm.Selected.Date.Day(), func(d dayInfo) int {
		return d.energy
	}, energyPalette))
	b.WriteString("\n")
	b.WriteString(renderLegend("Energy scale:", energyPalette))
	b.WriteString("\n")

	b.WriteString("Selected: ")
	b.WriteString(vm.Selected.Date.Format("2006-01-02"))
	b.WriteString("  mood ")
	b.WriteString(bar(vm.Selected.Mood))
	b.WriteString(fmt.Sprintf(" (%d)  energy ", vm.Selected.Mood))
	b.WriteString(bar(vm.Selected.Energy))
	b.WriteString(fmt.Sprintf(" (%d)\n", vm.Selected.Energy))

	return tea.NewView(b.String())
}
