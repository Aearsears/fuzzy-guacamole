package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#1a1a1a")).
			Bold(true).
			PaddingLeft(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87"))

	objectStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D9D9D9"))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3C3C3C")).
			Padding(0, 1).
			Margin(1, 2)
)

type model struct {
	buckets     []string
	selected    int
	viewObjects bool
	objects     []string
	s3Client    *s3.Client
	err         error
}

func initialModel(s3Client *s3.Client) model {
	// Load buckets on init
	ctx := context.Background()
	resp, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	var names []string
	if err == nil {
		for _, b := range resp.Buckets {
			names = append(names, *b.Name)
		}
	}
	return model{
		s3Client: s3Client,
		buckets:  names,
		selected: 0,
		err:      err,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.selected > 0 {
				m.selected--
			}

		case "down":
			if m.selected < len(m.buckets)-1 {
				m.selected++
			}

		case "enter":
			// Fetch objects for selected bucket
			bucket := m.buckets[m.selected]
			ctx := context.Background()
			resp, err := m.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
				Bucket:  aws.String(bucket),
				MaxKeys: 10,
			})
			if err != nil {
				m.err = err
				return m, nil
			}
			var objs []string
			for _, obj := range resp.Contents {
				objs = append(objs, fmt.Sprintf("%s (%d bytes)", *obj.Key, obj.Size))
			}
			m.viewObjects = true
			m.objects = objs
		case "backspace":
			m.viewObjects = false
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return borderStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	var b strings.Builder

	if !m.viewObjects {
		title := headerStyle.Render("S3 Buckets") + "  (↑ ↓ to navigate, Enter to open, q to quit)"
		b.WriteString(title + "\n\n")
		for i, name := range m.buckets {
			cursor := "  "
			style := objectStyle
			if i == m.selected {
				cursor = cursorStyle.Render("➜ ")
				style = selectedStyle
			}
			b.WriteString(style.Render(fmt.Sprintf("%s%s", cursor, name)) + "\n")
		}
	} else {
		title := headerStyle.Render(fmt.Sprintf("Objects in bucket: %s", m.buckets[m.selected])) +
			"  (Backspace to go back)"
		b.WriteString(title + "\n\n")
		if len(m.objects) == 0 {
			b.WriteString(objectStyle.Render("No objects found.\n"))
		} else {
			for _, obj := range m.objects {
				b.WriteString(objectStyle.Render("• "+obj) + "\n")
			}
		}
	}

	return borderStyle.Render(b.String())
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("Failed to load AWS config:", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	p := tea.NewProgram(initialModel(s3Client))
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
