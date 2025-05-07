package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/ini.v1"
)

type ProfileMenu struct {
	profiles        []string
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

	return ProfileMenu{
		profiles:        profiles,
		cursor:          0,
		selectedProfile: "",
	}
}

func (m ProfileMenu) Init() tea.Cmd {
	// todo: perform io loading in here
	return nil
}
func (m ProfileMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch {
		case key.Matches(msg, Keymap.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, Keymap.Down):
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}

		case key.Matches(msg, Keymap.Enter):
			m.selectedProfile = m.profiles[m.cursor]
			return m, func() tea.Msg {
				return ProfileMenuMessage{
					m.selectedProfile}
			}
		}
	}

	return m, nil
}
func (m ProfileMenu) View() string {
	// The header
	menu := fmt.Sprintf("Available AWS Profiles  \n\n")

	// Iterate over our choices
	for i, choice := range m.profiles {

		cursor := " " // no cursor
		display := ""

		if m.cursor == i {
			cursor = CursorStyle(">")              // cursor!
			display = SelectedStyle.Render(choice) // Highlight the selected choice
		} else {
			display = ChoiceStyle(choice) // Regular style for unselected choices
		}

		// Render the row with styles
		menu += fmt.Sprintf("%s %s\n", cursor, display)
	}

	// Send the UI for rendering
	return BorderStyle.Render(menu)
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
