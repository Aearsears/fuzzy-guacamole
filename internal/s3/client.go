package s3

import (
	"context"

	"github.com/Aearsears/fuzzy-guacamole/internal"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	tea "github.com/charmbracelet/bubbletea"
)

type Client struct {
	S3         *s3.Client
	NewMessage func() internal.APIMessage
	Wrapper    func(func() (any, error)) tea.Cmd
}

// client := Client{
// 	S3: nil, // your s3.Client here
// 	NewMessage: func() internal.APIMessage {
// 		return &S3MenuMessage{}
// 	},
// }

// client.Wrapper = func(fn func() (any,error)) tea.Cmd{
// 	return func() tea.Msg {
// 		output, err := fn()
// 		mssg := client.NewMessage()
// 		mssg.APIMessage = internal.APIMessage{
// 			Response: utils.FormatMetadata(output),
// 			Err:      err,
// 		}
// 		// format service specific message...
// 		return mssg

// 	}
// }

func (c *Client) PutObject(ctx context.Context, input *s3.PutObjectInput) tea.Cmd {
	return c.Wrapper(func() (any, error) {
		return c.S3.PutObject(ctx, input)
	})
}

func (c *Client) GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return c.S3.GetObject(ctx, input)
}

func (c *Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return c.S3.DeleteObject(ctx, input)
}
