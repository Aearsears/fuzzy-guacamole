package services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GoBackMessage struct{}

// services is a list of AWS services that are supported by the application
var services = []string{
	"S3",                // Object storage
	"Profiles",          // AWS credential profiles
	"EC2",               // Virtual servers
	"Lambda",            // Serverless functions
	"DynamoDB",          // NoSQL database
	"RDS",               // Relational databases
	"IAM",               // Identity management
	"CloudWatch",        // Monitoring and observability
	"VPC",               // Networking
	"API Gateway",       // API management
	"CloudFront",        // Content delivery
	"Route 53",          // DNS service
	"ECS",               // Container orchestration
	"EKS",               // Kubernetes service
	"SNS",               // Pub/sub messaging
	"SQS",               // Message queuing
	"CloudFormation",    // Infrastructure as code
	"Elastic Beanstalk", // App deployment
	"Secrets Manager",   // Secrets storage
	"CodePipeline",      // CI/CD pipeline
}

type MenuItem struct {
	name  string
	state SessionState
}

type MainMenu struct {
	choices         []MenuItem
	filteredChoices []MenuItem
	filterValue     string
	cursor          int
	selected        map[int]struct{}
	input           textinput.Model
}

func InitialMenu() MainMenu {
	menuItems := make([]MenuItem, len(services))
	for i, service := range services {
		menuItems[i] = MenuItem{
			name:  service,
			state: SessionState(i + 1),
		}
	}

	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = ""
	input.CharLimit = 250
	input.Width = 50

	return MainMenu{
		choices:         menuItems,
		filteredChoices: menuItems,
		cursor:          0,
		selected:        make(map[int]struct{}),
		input:           input,
		filterValue:     "",
	}
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.input.Focused() {
			// filter choices based on input
			if key.Matches(msg, Keymap.Enter) {
				m.filterValue = m.input.Value()

				if m.filterValue != "" {
					var filtered []MenuItem
					re, err := regexp.Compile("(?i)" + m.filterValue) // Case-insensitive
					if err == nil {
						for _, c := range m.choices {
							if re.MatchString(c.name) {
								filtered = append(filtered, c)
							}
						}
						m.filteredChoices = filtered
					}
					m.cursor = 0 // Reset cursor when filtering
				} else {
					m.filteredChoices = m.choices
				}
				m.input.SetValue("")
				m.input.Blur()
			}
			if key.Matches(msg, Keymap.Backspace) {
				m.filterValue = ""
				m.filteredChoices = m.choices
				m.input.SetValue("")
				m.input.Blur()
			}
			// only log keypresses for the input field when it's focused
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
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
				if len(m.filteredChoices) == 0 {
					return m, nil // No choices to select from

				}

				selected := m.filteredChoices[m.cursor]
				return m, func() tea.Msg {
					return SwitchMenuMessage{
						selected.state}
				}
			}
		}

		switch msg.String() {
		case "/":
			m.input.Focus()
			return m, textinput.Blink
		}

	}

	return m, tea.Batch(cmds...)
}

func (m MainMenu) View() string {
	// Get available width
	width := WindowSize.Width
	if width <= 0 {
		width = 80 // Default if window size not available
	}

	// Determine number of columns based on width
	numColumns := 3
	if width < 60 {
		numColumns = 1
	} else if width < 100 {
		numColumns = 2
	}

	// Calculate items per column
	totalItems := len(m.filteredChoices)
	itemsPerCol := (totalItems + numColumns - 1) / numColumns

	// Calculate column width
	contentWidth := width - 4 // Account for borders
	colWidth := contentWidth / numColumns

	// Build rows for the layout
	var rows []string

	for rowIdx := 0; rowIdx < itemsPerCol; rowIdx++ {
		var rowContent string

		for colIdx := 0; colIdx < numColumns; colIdx++ {
			// Calculate item index (reading down columns)
			itemIdx := rowIdx + colIdx*itemsPerCol

			if itemIdx < totalItems {
				choice := m.filteredChoices[itemIdx]

				// Style based on selection
				cursor := "  " // no cursor
				display := ""
				if m.cursor == itemIdx {
					cursor = CursorStyle("> ")
					display = SelectedStyle.Render(choice.name)
				} else {
					display = ChoiceStyle(choice.name)
				}

				// Add to row with fixed width
				item := fmt.Sprintf("%s%s", cursor, display)

				padding := colWidth - lipgloss.Width(item)
				if padding < 0 {
					padding = 0
				}
				paddedItem := item + strings.Repeat(" ", padding)
				rowContent += paddedItem

			}
		}

		rows = append(rows, rowContent)
	}

	menu := ""
	if len(rows) == 0 {
		rows = append(rows, DocStyle("No services available."))
	}
	menu = strings.Join(rows, "\n")

	if m.input.Focused() {
		menu += "\n" + m.input.View() // Add the input field at the bottom
	}
	return BorderStyle.Render(menu)
}
