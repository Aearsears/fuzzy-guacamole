package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/ini.v1"
)

type ProfileMenu struct {
	profiles        []string
	spinner         spinner.Model
	cursor          int
	selectedProfile string
}
type ProfileMenuMessage struct {
	profile string
}

func InitProfileMenu() ProfileMenu {
	profileSet := GetProfiles()
	profiles := make([]string, 0, len(profileSet))
	for key := range profileSet {
		profiles = append(profiles, key)
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return ProfileMenu{
		profiles:        profiles,
		cursor:          0,
		selectedProfile: "",
		spinner:         s,
	}
}

func (m ProfileMenu) Init() tea.Cmd {
	return m.spinner.Tick
}
func (m ProfileMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}

		case "enter":
			m.selectedProfile = m.profiles[m.cursor]
			return m, func() tea.Msg {
				return ProfileMenuMessage{
					m.selectedProfile}
			}
		}
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}
func (m ProfileMenu) View() string {
	// The header
	s := fmt.Sprintf("Available AWS Profiles %s \n\n", m.spinner.View())

	// Iterate over our choices
	for i, choice := range m.profiles {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if choice == m.selectedProfile {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// Send the UI for rendering
	return s
}

func getProfilesFromFile(path string, isConfig bool) ([]string, error) {
	profiles := []string{}
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	for _, section := range cfg.Sections() {
		name := section.Name()

		if name == "DEFAULT" {
			continue
		}

		// ~/.aws/config uses "profile foo"
		if isConfig {
			if strings.HasPrefix(name, "profile ") {
				profiles = append(profiles, strings.TrimPrefix(name, "profile "))
			}
		} else {
			// ~/.aws/credentials just uses "foo"
			profiles = append(profiles, name)
		}
	}

	return profiles, nil
}

func GetProfiles() map[string]struct{} {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".aws", "config")
	credsPath := filepath.Join(homeDir, ".aws", "credentials")

	configProfiles, err := getProfilesFromFile(configPath, true)
	if err != nil {
		log.Printf("error reading config file: %v", err)
	}

	credentialProfiles, err := getProfilesFromFile(credsPath, false)
	if err != nil {
		log.Printf("error reading credentials file: %v", err)
	}

	// Merge and deduplicate profiles
	profileSet := make(map[string]struct{})
	for _, p := range configProfiles {
		profileSet[p] = struct{}{}
	}
	for _, p := range credentialProfiles {
		profileSet[p] = struct{}{}
	}

	return profileSet
}
