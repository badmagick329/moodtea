package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func bar(n int) string {
	// n is 1..5
	filled := barFilledStyle.Render(strings.Repeat("█", n))
	empty := barEmptyStyle.Render(strings.Repeat("░", 5-n))
	return filled + empty
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

func renderDivider(width int) string {
	if width < 20 {
		width = 20
	}
	return dividerStyle.Render(strings.Repeat("─", width))
}

func renderTitleRow(title string, right string, width int) string {
	leftWidth := lipgloss.Width(title)
	rightWidth := lipgloss.Width(right)
	if width <= 0 {
		width = 60
	}
	gap := width - leftWidth - rightWidth
	if gap < 2 {
		return title + "\n" + subtleTextStyle.Render(right)
	}
	return title + strings.Repeat(" ", gap) + subtleTextStyle.Render(right)
}

func wrapChips(chips []string, width int) string {
	if width <= 0 {
		width = 60
	}
	var lines []string
	var current string
	currentWidth := 0
	for _, chip := range chips {
		w := lipgloss.Width(chip)
		if currentWidth > 0 && currentWidth+1+w > width {
			lines = append(lines, current)
			current = chip
			currentWidth = w
			continue
		}
		if currentWidth == 0 {
			current = chip
			currentWidth = w
			continue
		}
		current += " " + chip
		currentWidth += 1 + w
	}
	if current != "" {
		lines = append(lines, current)
	}
	return strings.Join(lines, "\n")
}

func renderStatChips(vm ViewModel, width int) string {
	chips := []string{
		chipStyle.Render(fmt.Sprintf("Avg M %.1f", vm.AvgMood)),
		chipStyle.Render(fmt.Sprintf("Avg E %.1f", vm.AvgEnergy)),
		chipStyle.Render(fmt.Sprintf("Med M %.1f", vm.MedianMood)),
		chipStyle.Render(fmt.Sprintf("Med E %.1f", vm.MedianEnergy)),
		chipStyle.Render(fmt.Sprintf("M %d-%d", vm.MinMood, vm.MaxMood)),
		chipStyle.Render(fmt.Sprintf("E %d-%d", vm.MinEnergy, vm.MaxEnergy)),
		chipStyle.Render(fmt.Sprintf("7d M %.1f", vm.Rolling7Mood)),
		chipStyle.Render(fmt.Sprintf("7d E %.1f", vm.Rolling7Energy)),
	}
	return wrapChips(chips, width)
}

func (m Model) renderHelp() string {
	if m.state.Mode == InputModeGotoMonth {
		return subtleTextStyle.Render(m.helpModel.View(gotoKeyMap{base: keys}))
	}
	return subtleTextStyle.Render(m.helpModel.View(keys))
}

func (m Model) renderHeaderCard(vm ViewModel) string {
	cardInnerWidth := m.state.Width - 4
	if cardInnerWidth < 50 {
		cardInnerWidth = 50
	}

	var b strings.Builder
	title := titleStyle.Render("MoodTea — " + vm.MonthLabel)
	recorded := fmt.Sprintf("%d days recorded", vm.RecordedDays)
	b.WriteString(renderTitleRow(title, recorded, cardInnerWidth))
	b.WriteString("\n")
	b.WriteString(renderStatChips(vm, cardInnerWidth))
	b.WriteString("\n")

	if m.state.Mode == InputModeGotoMonth {
		b.WriteString(subtleTextStyle.Render("Go to month (YYYY-MM): "))
		b.WriteString(m.state.GotoBuffer)
		if m.state.GotoError != "" {
			b.WriteString("  ")
			b.WriteString(warnTextStyle.Render(m.state.GotoError))
		}
		b.WriteString("\n")
	}

	b.WriteString(m.renderHelp())
	return headerCardStyle.Render(b.String())
}

func renderFooterCard(vm ViewModel) string {
	var b strings.Builder
	b.WriteString(subtleTextStyle.Render("Selected "))
	b.WriteString(vm.SelectedDate.Format("2006-01-02"))
	b.WriteString("\n")

	if vm.SelectedHasData {
		b.WriteString("Mood ")
		b.WriteString(bar(vm.Selected.Mood))
		b.WriteString(fmt.Sprintf(" (%d)  ", vm.Selected.Mood))
		b.WriteString("Energy ")
		b.WriteString(bar(vm.Selected.Energy))
		b.WriteString(fmt.Sprintf(" (%d)", vm.Selected.Energy))
	} else {
		b.WriteString(subtleTextStyle.Render("No entry for selected day"))
	}
	return footerCardStyle.Render(b.String())
}

func renderCalendarPanels(vm ViewModel, width int) string {
	mood := renderCalendar("Mood", vm.DaysInMonth, vm.StartWeekday, vm.DayMap, vm.SelectedDate.Day(), func(d dayInfo) int {
		return d.mood
	}, moodPalette) + "\n" + renderLegend("Mood scale:", moodPalette)

	energy := renderCalendar("Energy", vm.DaysInMonth, vm.StartWeekday, vm.DayMap, vm.SelectedDate.Day(), func(d dayInfo) int {
		return d.energy
	}, energyPalette) + "\n" + renderLegend("Energy scale:", energyPalette)

	const sideBySideMinWidth = 96
	if width >= sideBySideMinWidth {
		return lipgloss.JoinHorizontal(lipgloss.Top, mood, "    ", energy)
	}

	return mood + "\n\n" + energy
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

	b.WriteString(m.renderHeaderCard(vm))
	b.WriteString("\n")
	b.WriteString(renderDivider(m.state.Width))
	b.WriteString("\n\n\n")

	b.WriteString(renderCalendarPanels(vm, m.state.Width))
	b.WriteString("\n")
	b.WriteString(renderDivider(m.state.Width))
	b.WriteString("\n\n")
	b.WriteString(renderFooterCard(vm))
	b.WriteString("\n")

	m.viewport.SetContent(b.String())
	return tea.NewView(m.viewport.View())
}
