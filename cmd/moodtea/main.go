package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"moodtea/internal/data"
)

type model struct {
	days   []data.Day
	cursor int
	err    error
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "down", "j":
			if m.cursor < len(m.days)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func bar(n int) string {
	// n is 1..5
	return strings.Repeat("█", n) + strings.Repeat("░", 5-n)
}

func (m model) View() tea.View {
	if m.err != nil {
		return tea.NewView("Error: " + m.err.Error() + "\n\n(q to quit)\n")
	}
	if len(m.days) == 0 {
		return tea.NewView("No data.\n\n(q to quit)\n")
	}

	var b strings.Builder
	b.WriteString("MoodTea — January 2026\n")
	b.WriteString("↑/↓ (or k/j) to move, q to quit\n\n")

	// Render a simple “list”
	start := max(0, m.cursor-8)
	end := min(len(m.days), start+17)

	for i := start; i < end; i++ {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		d := m.days[i]
		b.WriteString(fmt.Sprintf(
			"%s%s  mood %s (%d)  energy %s (%d)\n",
			prefix,
			d.Date.Format("2006-01-02"),
			bar(d.Mood), d.Mood,
			bar(d.Energy), d.Energy,
		))
	}

	return tea.NewView(b.String())
}

func main() {
	var path string
	flag.StringVar(&path, "file", "data/january_2026.json", "path to JSON data file")
	flag.Parse()

	days, err := data.Load(path)
	m := model{days: days, err: err}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
