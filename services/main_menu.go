package services

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type GoBackMessage struct{}

// services is a list of AWS services that are supported by the application
var services = []string{"S3", "Profiles", "DynamoDB", "RDS", "Lambda", "SNS", "SQS", "CloudWatch", "IAM", "EC2"}

type MenuItem struct {
	name  string
	state SessionState
}

type MainMenu struct {
	choices  []MenuItem
	cursor   int
	selected map[int]struct{}
}

func InitialMenu() MainMenu {
	menuItems := make([]MenuItem, len(services))
	for i, service := range services {
		menuItems[i] = MenuItem{
			name:  service,
			state: SessionState(i + 1),
		}
	}
	return MainMenu{
		choices:  menuItems,
		cursor:   0,
		selected: make(map[int]struct{}),
	}
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, Keymap.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, Keymap.Down):
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case key.Matches(msg, Keymap.Enter):
			selected := m.choices[m.cursor]
			return m, func() tea.Msg {
				return SwitchMenuMessage{
					selected.state}
			}
		}

	}

	return m, nil
}

func (m MainMenu) View() string {
	menu := ""
	// Iterate over the menu choices
	for i, choice := range m.choices {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		display := ""

		if m.cursor == i {
			cursor = CursorStyle(">")                   // cursor!
			display = SelectedStyle.Render(choice.name) // Highlight the selected choice
		} else {
			display = ChoiceStyle(choice.name) // Regular style for unselected choices
		}

		// Render the row with styles
		menu += fmt.Sprintf("%s %s\n", cursor, display)
	}

	return menu

}
