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
	b.WriteString(fmt.Sprintf("Avg mood %.1f | Avg energy %.1f (%d days)\n", vm.AvgMood, vm.AvgEnergy, vm.RecordedDays))
	b.WriteString(fmt.Sprintf(
		"Median M/E: %.1f / %.1f | Min-Max M: %d-%d E: %d-%d | 7d Avg M/E: %.1f / %.1f\n",
		vm.MedianMood,
		vm.MedianEnergy,
		vm.MinMood,
		vm.MaxMood,
		vm.MinEnergy,
		vm.MaxEnergy,
		vm.Rolling7Mood,
		vm.Rolling7Energy,
	))
	if m.state.Mode == InputModeGotoMonth {
		b.WriteString(gotoHelpLine)
	} else {
		b.WriteString(helpLine)
	}
	b.WriteString("\n\n")

	b.WriteString(renderCalendar("Mood", vm.DaysInMonth, vm.StartWeekday, vm.DayMap, vm.SelectedDate.Day(), func(d dayInfo) int {
		return d.mood
	}, moodPalette))
	b.WriteString("\n")
	b.WriteString(renderLegend("Mood scale:", moodPalette))
	b.WriteString("\n\n")
	b.WriteString(renderCalendar("Energy", vm.DaysInMonth, vm.StartWeekday, vm.DayMap, vm.SelectedDate.Day(), func(d dayInfo) int {
		return d.energy
	}, energyPalette))
	b.WriteString("\n")
	b.WriteString(renderLegend("Energy scale:", energyPalette))
	b.WriteString("\n")

	if m.state.Mode == InputModeGotoMonth {
		b.WriteString("Go to month (YYYY-MM): ")
		b.WriteString(m.state.GotoBuffer)
		if m.state.GotoError != "" {
			b.WriteString("  ")
			b.WriteString(m.state.GotoError)
		}
		b.WriteString("\n")
	}

	if vm.SelectedHasData {
		b.WriteString("Selected: ")
		b.WriteString(vm.Selected.Date.Format("2006-01-02"))
		b.WriteString("  mood ")
		b.WriteString(bar(vm.Selected.Mood))
		b.WriteString(fmt.Sprintf(" (%d)  energy ", vm.Selected.Mood))
		b.WriteString(bar(vm.Selected.Energy))
		b.WriteString(fmt.Sprintf(" (%d)\n", vm.Selected.Energy))
	} else {
		b.WriteString("Selected: ")
		b.WriteString(vm.SelectedDate.Format("2006-01-02"))
		b.WriteString("  No entry for selected day\n")
	}

	return tea.NewView(b.String())
}
