package services

import (
	"github.com/charmbracelet/bubbles/key"
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

/* STYLING */

// DocStyle styling for viewports
var DocStyle = lipgloss.NewStyle().Margin(0, 2)

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
	PaddingLeft(1)

var ProfileStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("12")). // Red text
	Align(lipgloss.Right)

var ChoiceStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("7")) // Grey color for text

var CursorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("12")) // Red for the cursor

var SelectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("10")). // Green text for selected
	Italic(true)

var FooterStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8")) // Light grey footer

type keymap struct {
	Create key.Binding
	Enter  key.Binding
	Rename key.Binding
	Delete key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
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
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}
