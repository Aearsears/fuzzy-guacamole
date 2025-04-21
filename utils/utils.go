package utils

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Debug(msg string) {

	// fmt.Println(msg)
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	f.WriteString(msg + "\n")
}
