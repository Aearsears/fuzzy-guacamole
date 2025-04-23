package internal

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	tea "github.com/charmbracelet/bubbletea"
)

// wrapper for function to be called in a tea.cmd
func Wrapper(fn func() (*s3.PutObjectOutput, error)) tea.Cmd {

	return func() tea.Msg {
		res, err := fn()
		return MyMessage{
			buckets:     names,
			err:         err,
			loadBuckets: false}
	}
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
