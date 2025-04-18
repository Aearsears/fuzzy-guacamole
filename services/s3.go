package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styling constants for the S3 menu
var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5A5A5A")).
		//to fix: border goes beyon right edge
		MaxWidth(WindowSize.Width)

	leftPanel = lipgloss.NewStyle().
			Width(30).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5A5A5A")).
			MaxWidth(100)

	rightPanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5A5A5A")).
			MaxWidth(100)

	flexLayout = lipgloss.NewStyle().
			Align(lipgloss.Left)
)

type S3Menu struct {
	buckets     []string
	selected    int
	viewObjects bool
	objects     []string
	s3Client    *s3.Client
	err         error
	loading     bool
	spinner     spinner.Model
}

type S3MenuMessage struct {
	loadBuckets bool
	buckets     []string
	err         error
}

func loadBuckets(s3Client *s3.Client) tea.Cmd {

	ctx := context.Background()
	resp, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	var names []string
	if err == nil {
		for _, b := range resp.Buckets {
			names = append(names, *b.Name)
		}
	}

	return func() tea.Msg {
		return S3MenuMessage{
			buckets:     names,
			err:         err,
			loadBuckets: false,
		}
	}
}
func InitS3Menu() S3Menu {
	// Load buckets on init
	return S3Menu{
		s3Client: createS3Client(),
		buckets:  nil,
		objects:  nil,
		selected: 0,
		err:      nil,
		loading:  true,
		spinner:  CreateSpinner(),
	}
}

func (m S3Menu) Init() tea.Cmd {
	//todo : fix loading spinner
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		return S3MenuMessage{
			loadBuckets: true,
		}
	})
}

func (m S3Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {

	case S3MenuMessage:
		if msg.loadBuckets {
			cmds = append(cmds, loadBuckets(m.s3Client))
		}
		if msg.err != nil {
			m.err = msg.err
			m.loading = false
		} else if msg.buckets != nil {
			m.buckets = msg.buckets
			m.loading = false
		}

	case tea.KeyMsg:
		switch msg.String() {
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
				MaxKeys: aws.Int32(10),
			})
			// TODO: handle error in right side pane
			if err != nil {
				m.err = err
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
	return m, tea.Batch(cmds...)
}

func (m S3Menu) View() string {

	var left strings.Builder
	left.WriteString(HeaderStyle("Buckets") + "\n\n")
	if m.loading {
		left.WriteString(DocStyle(fmt.Sprintf("%s Loading buckets...\n", m.spinner.View())))
	} else if m.err != nil {
		left.WriteString(ErrStyle(fmt.Sprintf("Error: %v", m.err)))
	} else if len(m.buckets) == 0 {
		left.WriteString(DocStyle("No buckets found.\n"))
	} else {

		for i, name := range m.buckets {
			cursor := " "
			display := ""
			if i == m.selected {
				cursor = CursorStyle(">")
				display = SelectedStyle.Render(name)
			} else {
				display = ChoiceStyle(name)
			}
			left.WriteString(fmt.Sprintf("%s%s\n", cursor, display))
		}
	}

	var right strings.Builder
	if m.viewObjects {
		right.WriteString(HeaderStyle(fmt.Sprintf("Objects in: %s", m.buckets[m.selected])) + "\n\n")
		if len(m.objects) == 0 {
			right.WriteString(DocStyle("No objects found.\n"))
		} else {
			for _, obj := range m.objects {
				right.WriteString(DocStyle("• "+obj) + "\n")
			}
		}
	} else {
		right.WriteString(DocStyle("Press [Enter] to view bucket contents.\nUse ↑/↓ to navigate, q to quit."))
	}

	leftBox := leftPanel.Render(left.String())
	rightBox := rightPanel.Render(right.String())

	return flexLayout.Render(
		lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox),
	)
}

func createS3Client() *s3.Client {
	// Endpoint: http://localhost:4566
	// Region: us-east-1
	// Access key: test
	// Secret key: test
	// Custom endpoint resolver
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL:           "http://localhost:4566", // LocalStack endpoint
				SigningRegion: "us-east-1",
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})
	staticCreds := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider("test", "test", ""),
	)
	// Load config with custom resolver
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(staticCreds),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
		o.UsePathStyle = true
	})
	return s3Client
}
