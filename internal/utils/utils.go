package utils

import (
	"fmt"
	"os"

	"github.com/Aearsears/fuzzy-guacamole/internal"
	tea "github.com/charmbracelet/bubbletea"
)

// wrapper for function to be called in a tea.cmd
func Wrapper(fn func() (any, error)) tea.Cmd {
	return func() tea.Msg {
		output, err := fn()
		// maybe could do client.newMessage for each client's specific message type
		return internal.APIMessage{
			Response: FormatMetadata(output),
			Err:      err,
		}
	}
}

// todo: check types
func FormatMetadata(meta smithyhttp.ResponseMetadata) string {
	return fmt.Sprintf("Request ID: %s, Status: %d", meta.RequestID, meta.HTTPStatusCode)
}

func Debug(msg string) {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	f.WriteString(msg + "\n")
}
