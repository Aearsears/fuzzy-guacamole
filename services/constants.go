package services

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* CONSTANTS */

var (
	// P the current tea program
	P *tea.Program
	// WindowSize store the size of the terminal window
	WindowSize tea.WindowSizeMsg
)

func CreateSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

/* STYLING */

// DocStyle styling for viewports
var DocStyle = lipgloss.NewStyle().Margin(0, 2).Render

// HelpStyle styling for help context menu
var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

// ErrStyle provides styling for error messages
var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd534b")).Render

// AlertStyle provides styling for alert messages
var AlertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Render

var HeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7D56F4")). // Purple text
	Background(lipgloss.Color("#1a1a1a")). // Dark background
	Bold(true).
	PaddingLeft(1).
	Render

var ProfileStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("12")). // Red text
	Align(lipgloss.Right)

var ChoiceStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("7")). // Grey color for text
	Render

var CursorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("12")). // Red for the cursor
	Render

var SelectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("10")). // Green text for selected
	Italic(true)

var FooterStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8")). // Light grey footer
	Render

var BorderStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#5A5A5A")) // Grey border

// StatusBarStyle provides styling for the network activity status bar
var StatusBarStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#5A5A5A")). // Grey border
	Align(lipgloss.Left).
	Render

// StatusBarSuccessStyle indicates successful network operations
var StatusBarSuccessStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#0F4D0F")).
	Align(lipgloss.Left).
	Render

// StatusBarErrorStyle indicates failed network operations
var StatusBarErrorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#bd534b")).
	Align(lipgloss.Left).
	Render

type keymap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Create    key.Binding
	Enter     key.Binding
	Rename    key.Binding
	Delete    key.Binding
	Back      key.Binding
	Quit      key.Binding
	Backspace key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("up/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("down/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("left/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("right/l", "right"),
	),
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main menu"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Backspace: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "back to previous menu"),
	),
}
