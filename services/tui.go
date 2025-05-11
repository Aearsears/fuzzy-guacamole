package services

import (
	"fmt"

	"github.com/Aearsears/fuzzy-guacamole/internal"
	"github.com/Aearsears/fuzzy-guacamole/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type SessionState int

const (
	mainMenu SessionState = iota
	s3Menu
	profileMenu
)

type SwitchMenuMessage struct {
	menu SessionState
}

type TUI struct {
	state     SessionState
	views     map[SessionState]tea.Model
	profile   string
	config    aws.Config
	statusBar StatusBar
	// to implement
	quitting bool
}

func InitTUI() TUI {
	views := make(map[SessionState]tea.Model)
	views[mainMenu] = InitialMenu()
	cfg, _ := utils.LoadAWSConfig("")
	return TUI{
		state:     mainMenu,
		views:     views,
		profile:   "default",
		config:    cfg,
		quitting:  false,
		statusBar: InitStatusBar(),
	}
}

func (m TUI) Init() tea.Cmd {
	return nil
}

func (m TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case SwitchMenuMessage:
		if msg.menu == profileMenu {
			m.state = profileMenu
			if m.views[profileMenu] == nil {
				m.views[profileMenu] = InitProfileMenu()
				cmd = m.views[profileMenu].Init()
			}
		} else if msg.menu == s3Menu {
			m.state = s3Menu
			if m.views[s3Menu] == nil {
				m.views[s3Menu] = InitS3Menu()
				cmd = m.views[s3Menu].Init()
			}
		}
		return m, cmd
	case ProfileMenuMessage:
		m.state = mainMenu
		if m.profile != msg.profile {
			m.profile = msg.profile
			m.config = msg.config
			// if err != nil {
			// 	mssg := internal.APIMessage{
			// 		Err: err,
			// 	}
			// 	return m, func() tea.Msg {
			// 		return mssg
			// 	}
			// } else {
			// 	m.config = cfg
			// }

			return m, utils.SendMessage(internal.APIMessage{
				Status: fmt.Sprintf("Profile changed to %s", m.profile),
			})
		}
		// todo: if profile is different, then need to refresh all clients
		return m, nil

	case internal.APIMessage, StatusBarTimeoutMessage:
		newStatusBar, newCmd := m.statusBar.Update(msg)
		statusBar, ok := newStatusBar.(StatusBar)
		if !ok {
			panic("assertion on statusBar failed")
		}
		m.statusBar = statusBar
		return m, newCmd

	case tea.WindowSizeMsg:
		// todo: implement resize handling
		WindowSize = msg
		// top, right, bottom, left := DocStyle.GetMargin()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Quit):
			return m, tea.Quit

		case key.Matches(msg, Keymap.Back):
			//todo: implement go back to last screen, not mainmenu
			m.state = mainMenu
			return m, nil
		}
	}

	switch m.state {
	case mainMenu:
		newMainMenu, newCmd := m.views[mainMenu].Update(msg)
		mainMenuModel, ok := newMainMenu.(MainMenu)
		if !ok {
			panic("assertion on mainmenu failed")
		}
		m.views[mainMenu] = mainMenuModel
		cmd = newCmd
	case profileMenu:
		newProfile, newCmd := m.views[profileMenu].Update(msg)
		profileMenuModel, ok := newProfile.(ProfileMenu)
		if !ok {
			panic("assertion on profile menu failed")
		}
		m.views[profileMenu] = profileMenuModel
		cmd = newCmd
	case s3Menu:
		newS3, newCmd := m.views[s3Menu].Update(msg)
		s3MenuModel, ok := newS3.(S3Menu)
		if !ok {
			panic("assertion on S3 menu failed")
		}
		m.views[s3Menu] = s3MenuModel
		cmd = newCmd
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m TUI) View() string {
	menu := ""

	left := fmt.Sprintf("Profile: %s   Region: %s", m.profile, m.config.Region)
	center := ""
	switch m.state {
	case mainMenu:
		center = "[AWS] Main Menu"
	case profileMenu:
		center = "[AWS] Profiles"
	case s3Menu:
		center = "[AWS] S3"
	}

	totalWidth := WindowSize.Width
	leftWidth := lipgloss.Width(left)
	centerWidth := lipgloss.Width(center)

	// Calculate center start, but don't overlap left
	centerStart := (totalWidth - centerWidth) / 2
	if centerStart < leftWidth+1 {
		centerStart = leftWidth + 1
	}
	space := ""
	if centerStart > leftWidth {
		space = lipgloss.NewStyle().Width(centerStart - leftWidth).Render("")
	}

	headerLine := lipgloss.NewStyle().Width(leftWidth).Align(lipgloss.Left).Render(left) +
		space +
		lipgloss.NewStyle().Width(centerWidth).Align(lipgloss.Center).Bold(true).Render(center)

	menu += headerLine + "\n"

	switch m.state {
	case mainMenu:
		menu += m.views[mainMenu].View()
	case profileMenu:
		menu += m.views[profileMenu].View()
	case s3Menu:
		menu += m.views[s3Menu].View()
	}

	menu += "\n" + m.statusBar.View()

	helpText := "\n"
	for _, binding := range Keymap.List() {
		helpText += FooterStyle(fmt.Sprintf("%s ", binding.Help()))
	}

	menu += helpText
	return wordwrap.String(menu, WindowSize.Width)
}
