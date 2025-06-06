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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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
	buckets        []types.Bucket
	selected       int
	selectedBucket string
	viewObjects    bool
	objects        []string
	objectMetadata s3.S3ObjectMetadata
	paneFocus      int      // 0 = left for buckets, 1 = right for objects
	breadcrumbs    []string // stack of directories
	fileTree       *internal.Tree
	ptr            *internal.TreeNode
	savePath       string
	s3Client       s3.S3API
	err            error
	loading        bool
	spinner        spinner.Model
	input          textinput.Model
	createFlag     bool
}

func InitS3Menu() S3Menu {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Enter a new bucket name..."
	input.CharLimit = 250
	input.Width = 50

	//todo: handle error
	cfg, _ := utils.LoadAWSConfig("")

	client, _ := utils.ClientFactory("s3", cfg, true).(s3.S3API)
	return S3Menu{
		s3Client:    client,
		buckets:     nil,
		objects:     nil,
		fileTree:    &internal.Tree{},
		ptr:         &internal.TreeNode{},
		selected:    0,
		paneFocus:   0,
		breadcrumbs: []string{},
		err:         nil,
		loading:     true,
		spinner:     CreateSpinner(),
		input:       input,
		savePath:    ".",
		createFlag:  true,
	}
}

func (m S3Menu) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick,
		m.s3Client.ListBuckets(context.Background(),
			&s3aws.ListBucketsInput{}),
		utils.SendMessage(internal.APIMessage{
			Status: "Loading buckets...",
		}),
	)
}

// todo: allow user to change savePath
func (m S3Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {

	case internal.AWSConfigMessage:
		m.s3Client = m.createS3Client(msg.Config, true)
		//refresh the view
		// refresh last recently used views to not cause too much latency
		cmds = append(cmds,
			m.s3Client.ListBuckets(context.Background(),
				&s3aws.ListBucketsInput{}))

	case s3.S3MenuMessage:
		if msg.APIMessage.Err != nil {
			m.loading = false
			cmds = append(cmds, func() tea.Msg {
				return internal.APIMessage{
					Err: msg.APIMessage.Err,
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
			case s3.S3OpGetObjectMetadata:
				m.objectMetadata = msg.Metadata
				cmds = append(cmds, func() tea.Msg {
					return internal.APIMessage{
						Status: msg.APIMessage.Status,
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
				m.ptr = m.fileTree.Root
				m.paneFocus = 1
				m.selected = 0
				m.breadcrumbs = m.breadcrumbs[:0]
				m.breadcrumbs = append(m.breadcrumbs, m.ptr.Value)
				cmds = append(cmds, func() tea.Msg {
					return internal.APIMessage{
						Status: fmt.Sprintf("S3: Fetched %d objects successfully for %s", len(m.objects), m.selectedBucket),
					}
				})
			case s3.S3OpGetObject, s3.S3OpPutObject, s3.S3OpDeleteObject:
				cmds = append(cmds, func() tea.Msg {
					return internal.APIMessage{
						Status: msg.APIMessage.Status,
					}
				})

			}

		}

	case tea.KeyMsg:
		if m.input.Focused() {
			if key.Matches(msg, Keymap.Enter) {
				if m.paneFocus == 0 {
					cmds = append(cmds,
						m.s3Client.CreateBucket(context.Background(),
							&s3aws.CreateBucketInput{Bucket: aws.String(m.input.Value())}))
				} else {
					if m.createFlag {
						cmds = append(cmds,
							m.s3Client.PutObject(
								context.Background(),
								&s3aws.PutObjectInput{
									Bucket: aws.String(m.selectedBucket),
									Key:    aws.String(strings.Join(append(m.breadcrumbs[1:], m.input.Value()), "/"))},
								m.input.Value()))
					} else {
						if m.input.Value() == "y" {
							cmds = append(cmds,
								m.s3Client.DeleteObject(context.Background(),
									&s3aws.DeleteObjectInput{
										Bucket: aws.String(m.selectedBucket),
										Key:    aws.String(strings.Join(m.breadcrumbs[1:], "/")),
									}))
						}
						m.createFlag = true
					}
				}
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
			//bucket pane
			if m.paneFocus == 0 {
				switch {
				case key.Matches(msg, Keymap.Up):
					if m.selected > 0 {
						m.selected--
					}

				case key.Matches(msg, Keymap.Down):
					if m.selected < len(m.buckets)-1 {
						m.selected++
					}
				//todo: handle the case where the file tree is already visible
				case key.Matches(msg, Keymap.Left):
					if m.viewObjects {
						if m.selected > len(m.buckets)-1 {
							m.selected = len(m.buckets) - 1
						}

					}

				case key.Matches(msg, Keymap.Right):
					if m.viewObjects {
						m.paneFocus = 1
						//need to move select index to the last valid object in the array
						if m.selected > len(m.ptr.Children)-1 {
							m.selected = len(m.ptr.Children) - 1
						}
					}

				case key.Matches(msg, Keymap.Enter):
					if len(m.buckets) != 0 {
						m.selectedBucket = *m.buckets[m.selected].Name
						ctx := context.Background()
						cmds = append(cmds,
							m.s3Client.ListObjects(ctx,
								&s3aws.ListObjectsV2Input{Bucket: aws.String(m.selectedBucket),
									MaxKeys: aws.Int32(10)}))
					}

				case key.Matches(msg, Keymap.Create):
					m.input.Placeholder = "Enter a new bucket name..."
					m.input.Focus()
					cmds = append(cmds, textinput.Blink)

				case key.Matches(msg, Keymap.Backspace):
					m.viewObjects = false
					m.paneFocus = 0
				}
				switch msg.String() {
				case "r":
					cmds = append(cmds,
						m.s3Client.ListBuckets(context.Background(),
							&s3aws.ListBucketsInput{}))
				}
			} else {
				//object pane
				switch {
				case key.Matches(msg, Keymap.Up):
					if m.selected > 0 {
						m.selected--
					}

				case key.Matches(msg, Keymap.Down):
					if m.selected < len(m.ptr.Children)-1 {
						m.selected++
					}

				case key.Matches(msg, Keymap.Left):
					if m.ptr.Parent == nil {
						//at root, go back to buckets
						m.paneFocus = 0
						if m.selected > len(m.buckets)-1 || m.selected < 0 {
							m.selected = len(m.buckets) - 1
						}
					} else {
						//go up a level in the tree
						m.ptr = m.ptr.Parent
						m.breadcrumbs = m.breadcrumbs[:len(m.breadcrumbs)-1]
					}

				case key.Matches(msg, Keymap.Right):
					if len(m.ptr.Children) != 0 {
						//go down a level in the tree
						m.breadcrumbs = append(m.breadcrumbs, m.ptr.Children[m.selected].Value)
						m.ptr = m.ptr.Children[m.selected]
						m.selected = 0 // reset back to zero so dont get out of bounds
						if len(m.ptr.Children) == 0 {
							//get object metadata of file leaf node
							ctx := context.Background()
							cmds = append(cmds,
								m.s3Client.GetObjectMetadata(ctx,
									&s3aws.HeadObjectInput{
										Bucket: aws.String(m.selectedBucket),
										Key:    aws.String(strings.Join(m.breadcrumbs[1:], "/")),
									}))
						}
					}

				case key.Matches(msg, Keymap.Create):
					m.input.Placeholder = "Enter path of your file..."
					m.input.Focus()
					m.createFlag = true
					cmds = append(cmds, textinput.Blink)

				case key.Matches(msg, Keymap.Enter):
					if len(m.ptr.Children) == 0 && m.ptr.Value != "/" {
						ctx := context.Background()
						cmds = append(cmds,
							m.s3Client.GetObject(ctx,
								&s3aws.GetObjectInput{
									Bucket: aws.String(m.selectedBucket),
									Key:    aws.String(strings.Join(m.breadcrumbs[1:], "/")),
								}, m.savePath))
					}
					//TODO: somehow refresh the view after the file is downloaded/deleted
				case key.Matches(msg, Keymap.Delete):
					if len(m.ptr.Children) == 0 && m.ptr.Value != "/" {
						m.input.Placeholder = fmt.Sprintf("Confirm delete of %s [y/n]", strings.Join(m.breadcrumbs[1:], "/"))
						m.input.Focus()
						m.createFlag = false
						cmds = append(cmds, textinput.Blink)
					}

				}

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
		//todo create column for creation date and region
		for i, bucket := range m.buckets {
			cursor := " "
			display := ""

			if bucket.BucketRegion == nil {
				display = fmt.Sprintf("%s %s", *bucket.Name, bucket.CreationDate.Format("2006-01-02"))
			} else {
				display = fmt.Sprintf("%s %s %s", *bucket.Name, *bucket.BucketRegion, bucket.CreationDate.Format("2006-01-02"))
			}

			if i == m.selected && m.paneFocus == 0 {
				cursor = CursorStyle(">")
				display = SelectedStyle.Render(display)
			} else {
				display = ChoiceStyle(display)
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
			// render the current dir
			if len(m.ptr.Children) != 0 {
				for i, object := range m.ptr.Children {
					cursor := " "
					display := object.Value

					if i == m.selected && m.paneFocus == 1 {
						cursor = CursorStyle(">")
						display = SelectedStyle.Render(display)
					} else {
						display = ChoiceStyle(display)
					}

					right.WriteString(fmt.Sprintf("%s%s\n", cursor, display))
				}
			} else {
				//render the file metadata
				right.WriteString(fmt.Sprintf("File: %s\n", strings.Join(m.breadcrumbs[1:], "/")))
				right.WriteString(fmt.Sprintf("Size: %d bytes\n", m.objectMetadata.ContentLength))
				right.WriteString(fmt.Sprintf("Last Modified: %s\n", m.objectMetadata.LastModified.Format("2006-01-02 15:04:05")))
				right.WriteString(fmt.Sprintf("ETag: %s\n", m.objectMetadata.ETag))
				right.WriteString(fmt.Sprintf("Storage Class: %s\n", m.objectMetadata.StorageClass))
				right.WriteString(fmt.Sprintf("Content Type: %s\n", m.objectMetadata.ContentType))
				if len(m.objectMetadata.Metadata) != 0 {
					right.WriteString("User Metadata:\n")
					for k, v := range m.objectMetadata.Metadata {
						right.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
					}
				}
				right.WriteString(fmt.Sprintf("\nPress [Enter] to download %s\n", strings.Join(m.breadcrumbs[1:], "/")))
			}
		}
		right.WriteString("\n" + ChoiceStyle(m.breadcrumbs[0]+strings.Join(m.breadcrumbs[1:], "/")))
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

func (m S3Menu) createS3Client(cfg aws.Config, dev bool) s3.S3API {
	client, ok := utils.ClientFactory("s3", cfg, true).(s3.S3API)
	if !ok {
		panic("utils.ClientFactory(\"s3\") does not implement s3.S3API")
	}
	return client
}
