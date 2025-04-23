package services

import (
	"fmt"

	"github.com/Aearsears/fuzzy-guacamole/internal"
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
	views     map[int]tea.Model
	profile   string
	statusBar StatusBar
	// to implement
	quitting bool
}

func InitTUI() TUI {
	views := make(map[int]tea.Model)
	views[int(mainMenu)] = InitialMenu()
	return TUI{
		state:     mainMenu,
		views:     views,
		profile:   "N/A",
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
			if m.views[int(profileMenu)] == nil {
				m.views[int(profileMenu)] = InitProfileMenu()
				cmd = m.views[int(profileMenu)].Init()
			}
		} else if msg.menu == s3Menu {
			// switch to s3 menu and let it handle
			m.state = s3Menu
			if m.views[int(s3Menu)] == nil {
				m.views[int(s3Menu)] = InitS3Menu()
				cmd = m.views[int(s3Menu)].Init()
			}
		}
		return m, cmd
	case ProfileMenuMessage:
		m.profile = msg.profile
		// todo: if profile is different, then need to reautehntiate
		m.state = mainMenu
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
		newMainMenu, newCmd := m.views[int(mainMenu)].Update(msg)
		mainMenuModel, ok := newMainMenu.(MainMenu)
		if !ok {
			panic("assertion on mainmenu failed")
		}
		m.views[int(mainMenu)] = mainMenuModel
		cmd = newCmd
	case profileMenu:
		newProfile, newCmd := m.views[int(profileMenu)].Update(msg)
		profileMenuModel, ok := newProfile.(ProfileMenu)
		if !ok {
			panic("assertion on profile menu failed")
		}
		m.views[int(profileMenu)] = profileMenuModel
		cmd = newCmd
	case s3Menu:
		newS3, newCmd := m.views[int(s3Menu)].Update(msg)
		s3MenuModel, ok := newS3.(S3Menu)
		if !ok {
			panic("assertion on S3 menu failed")
		}
		m.views[int(s3Menu)] = s3MenuModel
		cmd = newCmd
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m TUI) View() string {
	menu := ""
	// todo; implement a better way to handle this
	switch m.state {
	case mainMenu:
		s := "[AWS] Main Menu"
		termWidth := lipgloss.Width(s) + 10           // Adding extra space to avoid clipping
		profileStyle := ProfileStyle.Width(termWidth) // Align text to the right
		menu += HeaderStyle(s) + " " + profileStyle.Render(fmt.Sprintf("Profile: %s", m.profile)) + "\n"
		menu += m.views[int(mainMenu)].View()
	case profileMenu:
		s := "[AWS] Profiles"
		termWidth := lipgloss.Width(s) + 10           // Adding extra space to avoid clipping
		profileStyle := ProfileStyle.Width(termWidth) // Align text to the right
		menu += HeaderStyle(s) + " " + profileStyle.Render(fmt.Sprintf("Profile: %s", m.profile)) + "\n"
		menu += m.views[int(profileMenu)].View()
	case s3Menu:
		s := "[AWS] S3"
		termWidth := lipgloss.Width(s) + 10           // Adding extra space to avoid clipping
		profileStyle := ProfileStyle.Width(termWidth) // Align text to the right
		menu += HeaderStyle(s) + " " + profileStyle.Render(fmt.Sprintf("Profile: %s", m.profile)) + "\n"
		menu += m.views[int(s3Menu)].View()
	}

	menu += "\n" + m.statusBar.View()

	helpText := "\n"
	helpText += FooterStyle(fmt.Sprintf("%s ", Keymap.Up.Help()))
	helpText += FooterStyle(fmt.Sprintf("%s ", Keymap.Down.Help()))
	helpText += FooterStyle(fmt.Sprintf("%s ", Keymap.Enter.Help()))
	helpText += FooterStyle(fmt.Sprintf("%s ", Keymap.Back.Help()))
	helpText += FooterStyle(fmt.Sprintf("%s ", Keymap.Backspace.Help()))
	helpText += FooterStyle(fmt.Sprintf("%s \n", Keymap.Quit.Help()))

	menu += helpText
	return wordwrap.String(menu, WindowSize.Width)
}
