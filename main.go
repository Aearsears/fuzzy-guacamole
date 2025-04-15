package main

import (
	"fmt"
	"os"

	"github.com/Aearsears/fuzzy-guacamole/services"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(services.InitialMenu())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
