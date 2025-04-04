package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
	tea "github.com/charmbracelet/bubbletea"
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
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	var b strings.Builder

	if !m.viewObjects {
		b.WriteString("S3 Buckets (use ↑ ↓ to navigate, Enter to open, q to quit):\n\n")
		for i, name := range m.buckets {
			cursor := "  "
			if i == m.selected {
				cursor = "➜ "
			}
			b.WriteString(fmt.Sprintf("%s%s\n", cursor, name))
		}
	} else {
		b.WriteString(fmt.Sprintf("Objects in bucket: %s (press Backspace to go back)\n\n", m.buckets[m.selected]))
		if len(m.objects) == 0 {
			b.WriteString("No objects found.\n")
		} else {
			for _, obj := range m.objects {
				b.WriteString(fmt.Sprintf("• %s\n", obj))
			}
		}
	}

	return b.String()
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
