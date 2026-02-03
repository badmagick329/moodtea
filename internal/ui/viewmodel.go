package ui

import "moodtea/internal/data"

type ViewModel struct {
	MonthLabel   string
	DaysInMonth  int
	StartWeekday int
	DayMap       map[int]dayInfo
	Selected     data.Day
	HasData      bool
}

func (m Model) buildViewModel() ViewModel {
	days := m.currentDays()
	if len(days) == 0 {
		return ViewModel{HasData: false}
	}

	start, daysInMonth, startWeekday := monthBounds(days)
	dayMap := buildDayMap(days)
	selected := days[m.state.Cursor]

	return ViewModel{
		MonthLabel:   start.Format("January 2006"),
		DaysInMonth:  daysInMonth,
		StartWeekday: startWeekday,
		DayMap:       dayMap,
		Selected:     selected,
		HasData:      true,
	}
}
