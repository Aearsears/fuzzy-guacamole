package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/Aearsears/fuzzy-guacamole/internal"
	"github.com/Aearsears/fuzzy-guacamole/internal/s3"
	"github.com/Aearsears/fuzzy-guacamole/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	s3aws "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmbracelet/bubbles/key"
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
	buckets        []string
	selected       int
	selectedBucket string
	viewObjects    bool
	objects        []string
	fileTree       *internal.Tree
	s3Client       s3.S3API
	err            error
	loading        bool
	spinner        spinner.Model
	input          textinput.Model
}

func InitS3Menu() S3Menu {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Enter a new bucket name..."
	input.CharLimit = 250
	input.Width = 50
	return S3Menu{
		s3Client: createS3Client(),
		buckets:  nil,
		objects:  nil,
		fileTree: &internal.Tree{},
		selected: 0,
		err:      nil,
		loading:  true,
		spinner:  CreateSpinner(),
		input:    input,
	}
}

func (m S3Menu) Init() tea.Cmd {
	// Load buckets on init
	return tea.Batch(m.spinner.Tick,
		m.s3Client.ListBuckets(context.Background(),
			&s3aws.ListBucketsInput{}),
		utils.SendMessage(internal.APIMessage{
			Status: "Loading buckets...",
		}),
	)
}

// todo: get, put and delete objects
func (m S3Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {

	case s3.S3MenuMessage:
		if msg.APIMessage.Err != nil {
			m.err = msg.APIMessage.Err
			m.loading = false
			cmds = append(cmds, func() tea.Msg {
				return internal.APIMessage{
					Err: m.err,
				}
			})

		} else {
			m.loading = false
			switch msg.Op {
			case s3.S3OpListBuckets:
				m.buckets = msg.Buckets
				cmds = append(cmds, func() tea.Msg {
					return internal.APIMessage{
						Status: fmt.Sprintf("S3: Listed %d buckets successfully", len(m.buckets)),
					}
				})
			case s3.S3OpCreateBucket:
				cmds = append(cmds, func() tea.Msg {
					return internal.APIMessage{
						Status: fmt.Sprintf("S3: Created bucket %s", msg.Bucket),
					}
				}, m.s3Client.ListBuckets(context.Background(),
					&s3aws.ListBucketsInput{}))
			case s3.S3OpListObjects:
				m.viewObjects = true
				m.objects = msg.Objects
				m.fileTree = internal.CreateTree(m.objects)
				cmds = append(cmds, func() tea.Msg {
					return internal.APIMessage{
						Status: fmt.Sprintf("S3: Fetched %d objects successfully for %s", len(m.objects), m.selectedBucket),
					}
				})
			}

		}

	case tea.KeyMsg:
		if m.input.Focused() {
			//todo: handle multiple operations for buckets in the same input field
			if key.Matches(msg, Keymap.Enter) {
				cmds = append(cmds,
					m.s3Client.CreateBucket(context.Background(),
						&s3aws.CreateBucketInput{Bucket: aws.String(m.input.Value())}))
				m.input.SetValue("")
				m.input.Blur()
			}
			if key.Matches(msg, Keymap.Backspace) {
				m.input.SetValue("")
				m.input.Blur()
			}
			// only log keypresses for the input field when it's focused
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, Keymap.Up):
				if m.selected > 0 {
					m.selected--
				}

			case key.Matches(msg, Keymap.Down):
				if m.selected < len(m.buckets)-1 {
					m.selected++
				}

			case key.Matches(msg, Keymap.Enter):
				if len(m.buckets) != 0 {
					m.selectedBucket = m.buckets[m.selected]
					ctx := context.Background()
					cmds = append(cmds,
						m.s3Client.ListObjects(ctx,
							&s3aws.ListObjectsV2Input{Bucket: aws.String(m.selectedBucket),
								MaxKeys: aws.Int32(10)}))
				}

			case key.Matches(msg, Keymap.Create):
				m.input.Focus()
				cmds = append(cmds, textinput.Blink)

			case key.Matches(msg, Keymap.Backspace):
				m.viewObjects = false
			}
			switch msg.String() {
			case "r":
				cmds = append(cmds,
					m.s3Client.ListBuckets(context.Background(),
						&s3aws.ListBucketsInput{}))
			}

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
		left.WriteString(ErrStyle("Error: :c"))
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
	// would be cool if could view objects like a tree from left to right
	if m.viewObjects {
		right.WriteString(HeaderStyle(fmt.Sprintf("Objects in: %s", m.selectedBucket)) + "\n\n")
		if len(m.objects) == 0 {
			right.WriteString(DocStyle("No objects found.\n"))
		} else {
			right.WriteString(DocStyle(m.fileTree.Display()))
		}
	} else {
		right.WriteString(DocStyle("Press [Enter] to view bucket contents."))
	}

	leftBox := leftPanel.Render(left.String())
	rightBox := rightPanel.Render(right.String())
	menu := flexLayout.Render(
		lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox),
	)

	if m.input.Focused() {
		menu += "\n" + m.input.View() // Add the input field at the bottom
	}
	return menu
}

func createS3Client() s3.S3API {
	client, ok := utils.ClientFactory("s3").(s3.S3API)
	if !ok {
		panic("utils.ClientFactory(\"s3\") does not implement s3.S3API")
	}
	return client
}
