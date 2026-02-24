package ui

type InputMode int

const (
	InputModeNormal InputMode = iota
	InputModeGotoMonth
)

type State struct {
	MonthIndex int
	CursorDay  int
	Mode       InputMode
	GotoBuffer string
	GotoError  string
	Width      int
	Height     int
}
