package services

import (
	"time"

	"github.com/Aearsears/fuzzy-guacamole/internal"
	tea "github.com/charmbracelet/bubbletea"
)

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
		timeout: 5,
		loading: false,
	}
}

func (m StatusBar) Init() tea.Cmd {
	return nil
}

func (m StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case internal.APIMessage:
		// todo: how to timeout multiple messsages at a time?
		m.display = true
		if msg.Status != "" {
			m.display_text = msg.Status
			m.loading = true
		} else if msg.Err != nil {
			m.err = msg.Err
			m.display_text = m.err.Error()
			m.loading = false
		} else if msg.Response != "" {
			if response, ok := msg.Response.(string); ok {
				m.display_text = response
			} else {
				m.display_text = "Invalid response type"
			}
			m.loading = false
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
		if m.loading {
			return StatusBarStyle(m.display_text)
		} else if m.err != nil {
			return StatusBarErrorStyle(m.err.Error())
		} else {
			return StatusBarSuccessStyle(m.display_text)
		}
	}
	return ""
}
