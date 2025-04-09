package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainMenu struct {
	state    string
	choices  []string
	cursor   int
	selected map[int]struct{}
	profile  string
	submenu  *ProfileMenu
}

func initialMenu() MainMenu {
	return MainMenu{
		state:    "main",
		choices:  []string{"S3", "DynamoDb", "Profiles"},
		cursor:   0,
		selected: make(map[int]struct{}),
		profile:  "",
	}
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		}
		switch m.state {
		case "main":
			switch msg.String() {
			case "enter", " ":
				menu := m.choices[m.cursor]
				if menu == "Profiles" {
					// Switch to submenu and let it handle
					m.state = "submenu"
					submenu := InitProfileMenu()
					m.submenu = &submenu
					return m.submenu, nil
				}
			}
		case "submenu":
			// After returning from submenu, get the returned value (e.g., Profile)
			if m.submenu != nil && m.submenu.selectedProfile != "" {
				// Set the returned value in the main menu's profileValue
				m.profile = m.submenu.selectedProfile
				// Return to the main menu
				m.state = "main"
				m.submenu = nil
				return m, nil
			}
		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m MainMenu) profileMenu() {
	p := tea.NewProgram(InitProfileMenu())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m MainMenu) View() string {
	// Define the header style
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")). // Purple text
		Background(lipgloss.Color("#1a1a1a")). // Dark background
		Bold(true).
		PaddingLeft(1)
	// Get the terminal width to align "orofile" to the far right
	termWidth := lipgloss.Width(headerStyle.Render("[AWS] Main Menu")) + 10 // Adding extra space to avoid clipping

	profileStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")). // Red text for "orofile"
		Align(lipgloss.Right).
		Width(termWidth) // Align text to the right
	// Define the choice style (for regular menu items)
	choiceStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")) // Grey color for text

	// Define the cursor style (for selected menu item)
	cursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")) // Red for the cursor

	// Define the selected item style
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")). // Green text for selected
		Italic(true)

		// The header with aligned "orofile"
	s := headerStyle.Render("[AWS] Main Menu") + " " + profileStyle.Render(fmt.Sprintf("Profile: %s", m.profile)) + "\n\n"

	// Iterate over the menu choices
	for i, choice := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = cursorStyle.Render(">") // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = selectedStyle.Render("x") // selected!
		}

		// Render the row with styles
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choiceStyle.Render(choice))
	}

	// Footer with styling
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // Light grey footer
	s += "\n" + footerStyle.Render("Press q to quit.\n")

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialMenu())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
