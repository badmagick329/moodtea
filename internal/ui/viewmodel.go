package ui

import "moodtea/internal/data"

type ViewModel struct {
	MonthLabel   string
	DaysInMonth  int
	StartWeekday int
	DayMap       map[int]dayInfo
	Selected     data.Day
	AvgMood      float64
	AvgEnergy    float64
	RecordedDays int
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
	avgMood, avgEnergy, recorded := monthAverages(days)

	return ViewModel{
		MonthLabel:   start.Format("January 2006"),
		DaysInMonth:  daysInMonth,
		StartWeekday: startWeekday,
		DayMap:       dayMap,
		Selected:     selected,
		AvgMood:      avgMood,
		AvgEnergy:    avgEnergy,
		RecordedDays: recorded,
		HasData:      true,
	}
}

func monthAverages(days []data.Day) (avgMood float64, avgEnergy float64, recorded int) {
	recorded = len(days)
	if recorded == 0 {
		return 0, 0, 0
	}

	var moodSum int
	var energySum int
	for _, day := range days {
		moodSum += day.Mood
		energySum += day.Energy
	}

	avgMood = float64(moodSum) / float64(recorded)
	avgEnergy = float64(energySum) / float64(recorded)
	return avgMood, avgEnergy, recorded
}
