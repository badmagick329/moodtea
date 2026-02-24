package ui

import (
	"time"

	"moodtea/internal/data"
)

type ViewModel struct {
	MonthLabel      string
	DaysInMonth     int
	StartWeekday    int
	DayMap          map[int]dayInfo
	SelectedDate    time.Time
	Selected        data.Day
	SelectedHasData bool
	AvgMood         float64
	AvgEnergy       float64
	MedianMood      float64
	MedianEnergy    float64
	MinMood         int
	MaxMood         int
	MinEnergy       int
	MaxEnergy       int
	Rolling7Mood    float64
	Rolling7Energy  float64
	RecordedDays    int
	HasData         bool
}

func (m Model) buildViewModel() ViewModel {
	days := m.currentDays()
	if len(days) == 0 {
		return ViewModel{HasData: false}
	}

	start, daysInMonth, startWeekday := monthBounds(days)
	dayMap := buildDayMap(days)
	selectedDate := m.currentDateForCursor()
	selected, selectedHasData := m.currentSelectedEntry()
	avgMood, avgEnergy, recorded := monthAverages(days)
	medianMood, medianEnergy := monthMedians(days)
	minMood, maxMood, minEnergy, maxEnergy := monthMinMax(days)
	rolling7Mood, rolling7Energy := monthRolling7(days)

	return ViewModel{
		MonthLabel:      start.Format("January 2006"),
		DaysInMonth:     daysInMonth,
		StartWeekday:    startWeekday,
		DayMap:          dayMap,
		SelectedDate:    selectedDate,
		Selected:        selected,
		SelectedHasData: selectedHasData,
		AvgMood:         avgMood,
		AvgEnergy:       avgEnergy,
		MedianMood:      medianMood,
		MedianEnergy:    medianEnergy,
		MinMood:         minMood,
		MaxMood:         maxMood,
		MinEnergy:       minEnergy,
		MaxEnergy:       maxEnergy,
		Rolling7Mood:    rolling7Mood,
		Rolling7Energy:  rolling7Energy,
		RecordedDays:    recorded,
		HasData:         true,
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

func monthMedians(days []data.Day) (float64, float64) {
	if len(days) == 0 {
		return 0, 0
	}

	moods := make([]int, 0, len(days))
	energies := make([]int, 0, len(days))
	for _, day := range days {
		moods = append(moods, day.Mood)
		energies = append(energies, day.Energy)
	}
	return median(moods), median(energies)
}

func monthMinMax(days []data.Day) (int, int, int, int) {
	if len(days) == 0 {
		return 0, 0, 0, 0
	}
	minMood, maxMood := days[0].Mood, days[0].Mood
	minEnergy, maxEnergy := days[0].Energy, days[0].Energy
	for _, day := range days[1:] {
		if day.Mood < minMood {
			minMood = day.Mood
		}
		if day.Mood > maxMood {
			maxMood = day.Mood
		}
		if day.Energy < minEnergy {
			minEnergy = day.Energy
		}
		if day.Energy > maxEnergy {
			maxEnergy = day.Energy
		}
	}
	return minMood, maxMood, minEnergy, maxEnergy
}

func monthRolling7(days []data.Day) (float64, float64) {
	if len(days) == 0 {
		return 0, 0
	}
	start := len(days) - 7
	if start < 0 {
		start = 0
	}
	window := days[start:]
	avgMood, avgEnergy, _ := monthAverages(window)
	return avgMood, avgEnergy
}

func median(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]int(nil), values...)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j] < sorted[i] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return float64(sorted[mid])
	}
	return float64(sorted[mid-1]+sorted[mid]) / 2
}
