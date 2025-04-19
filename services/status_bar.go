package services

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type APIMessage struct {
	err      error
	response string
}
type StatusBarTimeoutMessage struct{}

type StatusBar struct {
	err          error
	display_text string
	display      bool
	loading      bool
	timeout      int
}

func statusBarTimeout(seconds int) tea.Cmd {
	return tea.Tick(time.Duration(seconds)*time.Second, func(time.Time) tea.Msg {
		return StatusBarTimeoutMessage{}
	})
}

func InitStatusBar() StatusBar {
	return StatusBar{
		timeout: 10,
		loading: false,
	}
}

func (m StatusBar) Init() tea.Cmd {
	return nil
}

func (m StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case APIMessage:
		m.display = true
		if msg.err != nil {
			m.err = msg.err
			m.display_text = m.err.Error()
		} else {
			m.display_text = msg.response
		}
		return m, statusBarTimeout(m.timeout)

	case StatusBarTimeoutMessage:
		m.display = false
		return m, nil
	}
	return m, nil
}

func (m StatusBar) View() string {
	if m.display {
		if m.err != nil {
			return StatusBarErrorStyle(m.err.Error())
		} else {
			return StatusBarSuccessStyle(m.display_text)
		}
	}
	return ""
}
