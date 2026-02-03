package ui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"moodtea/internal/data"
)

type dayInfo struct {
	mood   int
	energy int
	has    bool
}

func buildDayMap(days []data.Day) map[int]dayInfo {
	out := make(map[int]dayInfo, len(days))
	for _, d := range days {
		out[d.Date.Day()] = dayInfo{
			mood:   d.Mood,
			energy: d.Energy,
			has:    true,
		}
	}
	return out
}

func monthBounds(days []data.Day) (time.Time, int, int) {
	first := days[0].Date
	loc := first.Location()
	start := time.Date(first.Year(), first.Month(), 1, 0, 0, 0, 0, loc)
	daysInMonth := time.Date(first.Year(), first.Month()+1, 0, 0, 0, 0, 0, loc).Day()
	startWeekday := int(start.Weekday()) // Sunday=0
	return start, daysInMonth, startWeekday
}

func renderCalendar(title string, daysInMonth int, startWeekday int, dayMap map[int]dayInfo, selectedDay int, value func(dayInfo) int, palette []string) string {
	const cellWidth = 4
	weekdayHeader := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	headerStyle := lipgloss.NewStyle().Bold(true)
	titleStyle := lipgloss.NewStyle().Bold(true)
	cellStyle := lipgloss.NewStyle().Width(cellWidth).Align(lipgloss.Center)
	emptyStyle := lipgloss.NewStyle().Width(cellWidth)

	var b strings.Builder
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n")
	for _, w := range weekdayHeader {
		b.WriteString(headerStyle.Render(cellStyle.Render(w)))
	}
	b.WriteString("\n")

	totalCells := 6 * 7
	for i := 0; i < totalCells; i++ {
		day := i - startWeekday + 1
		if day < 1 || day > daysInMonth {
			b.WriteString(emptyStyle.Render(""))
		} else {
			info := dayMap[day]
			style := cellStyle
			if info.has {
				val := value(info)
				if val >= 1 && val <= len(palette) {
					style = style.Background(lipgloss.Color(palette[val-1])).Foreground(lipgloss.Color("#FFFFFF"))
				}
			}
			if day == selectedDay {
				style = style.Bold(true).Underline(true)
			}
			b.WriteString(style.Render(fmt.Sprintf("%2d", day)))
		}
		if (i+1)%7 == 0 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
