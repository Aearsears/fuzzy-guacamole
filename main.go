package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SessionState int

const (
	mainMenu SessionState = iota
	profileMenu
)

type MainMenu struct {
	state    SessionState
	views    map[int]tea.Model
	choices  []string
	cursor   int
	selected map[int]struct{}
	profile  string
	// to implement
	quitting bool
}

func initialMenu() MainMenu {
	return MainMenu{
		state:    mainMenu,
		views:    make(map[int]tea.Model),
		choices:  []string{"S3", "DynamoDb", "Profiles"},
		cursor:   0,
		selected: make(map[int]struct{}),
		profile:  "",
		quitting: false,
	}
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ProfileMenuMessage:
		m.profile = msg.profile
		// todo: if profile is different, then need to reautehntiate
		m.state = mainMenu
		return m, nil
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	switch m.state {
	case mainMenu:
		if m.views[int(mainMenu)] == nil {
			m.views[int(mainMenu)] = m
		}

		switch msg := msg.(type) {

		case tea.KeyMsg:

			switch msg.String() {

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case "enter":
				selected := m.choices[m.cursor]
				if selected == "Profiles" {
					// Switch to submenu and let it handle
					m.state = profileMenu
					if m.views[int(profileMenu)] == nil {
						m.views[int(profileMenu)] = InitProfileMenu()
						cmd = m.views[int(profileMenu)].Init()
					}
					return m, cmd
				}
			}

		}
	case profileMenu:
		newProfile, newCmd := m.views[int(profileMenu)].Update(msg)
		profileMenuModel, ok := newProfile.(ProfileMenu)
		if !ok {
			panic("assertion on profile menu failed")
		}
		m.views[int(profileMenu)] = profileMenuModel
		cmd = newCmd
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MainMenu) View() string {
	switch m.state {
	case profileMenu:
		return m.views[int(profileMenu)].View()
	default:
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
}

func main() {
	p := tea.NewProgram(initialMenu())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
