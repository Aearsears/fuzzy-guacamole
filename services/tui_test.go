package services

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// Mock implementations for menu models
type mockModel struct{ tea.Model }

func (m mockModel) Init() tea.Cmd                           { return nil }
func (m mockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m mockModel) View() string                            { return "mock view" }

func mockStatusBar() StatusBar {
	return StatusBar{}
}

func TestInitTUI(t *testing.T) {
	// Override menu initializers for test
	InitialMenu := func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	assert.Equal(t, mainMenu, tui.state)
	assert.Equal(t, "default", tui.profile)
	assert.NotNil(t, tui.views[mainMenu])
}

func TestUpdateSwitchMenuMessage(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	msg := SwitchMenuMessage{menu: profileMenu}
	updated, _ := tui.Update(msg)
	tui2 := updated.(TUI)
	assert.Equal(t, profileMenu, tui2.state)
	assert.NotNil(t, tui2.views[profileMenu])
}

func TestUpdateProfileMenuMessage(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	msg := ProfileMenuMessage{
		profile: "test-profile",
		config:  aws.Config{Region: "us-west-2"},
	}
	updated, _ := tui.Update(msg)
	tui2 := updated.(TUI)
	assert.Equal(t, mainMenu, tui2.state)
	assert.Equal(t, "test-profile", tui2.profile)
	assert.Equal(t, "us-west-2", tui2.config.Region)
}

func TestViewRenders(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	WindowSize.Width = 80
	view := tui.View()
	assert.Contains(t, view, "Profile: default")
	assert.Contains(t, view, "[AWS]")
}

// Additional tests

func TestUpdateS3MenuSwitch(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	msg := SwitchMenuMessage{menu: s3Menu}
	updated, _ := tui.Update(msg)
	tui2 := updated.(TUI)
	assert.Equal(t, s3Menu, tui2.state)
	assert.NotNil(t, tui2.views[s3Menu])
}

func TestUpdateKeyMsgQuit(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	// Simulate a quit key press
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	Keymap.Quit = key.NewBinding(key.WithKeys("ctrl+c"))
	model, cmd := tui.Update(msg)
	assert.Equal(t, tui, model)
	assert.NotNil(t, cmd)
}

func TestUpdateKeyMsgBack(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	tui.state = profileMenu
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	Keymap.Back = key.NewBinding(key.WithKeys("esc"))
	model, cmd := tui.Update(msg)
	tui2 := model.(TUI)
	assert.Equal(t, mainMenu, tui2.state)
	assert.Nil(t, cmd)
}

func TestUpdateStatusBarTimeoutMessage(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	msg := StatusBarTimeoutMessage{}
	model, _ := tui.Update(msg)
	_, ok := model.(TUI)
	assert.True(t, ok)
}

func TestUpdateWindowSizeMsg(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	msg := tea.WindowSizeMsg{Width: 100, Height: 40}
	model, _ := tui.Update(msg)
	_, ok := model.(TUI)
	assert.True(t, ok)
}

func TestUpdateWithUnknownMsg(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	type unknownMsg struct{}
	model, _ := tui.Update(unknownMsg{})
	_, ok := model.(TUI)
	assert.True(t, ok)
}

func TestUpdateSwitchMenuMessage(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	msg := SwitchMenuMessage{menu: profileMenu}
	updated, _ := tui.Update(msg)
	tui2 := updated.(TUI)
	assert.Equal(t, profileMenu, tui2.state)
	assert.NotNil(t, tui2.views[profileMenu])
}

func TestUpdateProfileMenuMessage(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	msg := ProfileMenuMessage{
		profile: "test-profile",
		config:  aws.Config{Region: "us-west-2"},
	}
	updated, _ := tui.Update(msg)
	tui2 := updated.(TUI)
	assert.Equal(t, mainMenu, tui2.state)
	assert.Equal(t, "test-profile", tui2.profile)
	assert.Equal(t, "us-west-2", tui2.config.Region)
}

func TestViewRenders(t *testing.T) {
	InitialMenu = func() tea.Model { return mockModel{} }
	InitProfileMenu = func() tea.Model { return mockModel{} }
	InitS3Menu = func() tea.Model { return mockModel{} }
	InitStatusBar = mockStatusBar

	tui := InitTUI()
	WindowSize.Width = 80
	view := tui.View()
	assert.Contains(t, view, "Profile: default")
	assert.Contains(t, view, "[AWS]")
}
