package services

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
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
	totalItems := len(m.choices)
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
				choice := m.choices[itemIdx]

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

	menu := strings.Join(rows, "\n")
	return BorderStyle.Render(menu)
}
