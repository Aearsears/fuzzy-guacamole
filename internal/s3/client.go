package s3

import (
	"context"

	"github.com/Aearsears/fuzzy-guacamole/internal"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	tea "github.com/charmbracelet/bubbletea"
)

type S3Client struct {
	Client *s3.Client
}

type S3MenuMessage struct {
	APIMessage   internal.APIMessage
	CreateBucket bool
	LoadBuckets  bool
	Buckets      []string
}

func (c *S3Client) NewMessage() S3MenuMessage {
	return S3MenuMessage{
		APIMessage:   internal.APIMessage{},
		CreateBucket: false,
		LoadBuckets:  false,
		Buckets:      nil,
	}
}

func (c *S3Client) Wrapper(fn func() (any, error)) tea.Cmd {
	return func() tea.Msg {
		a, _ := fn()
		return a
	}
}

func (c *S3Client) ListBuckets(ctx context.Context, input *s3.ListBucketsInput) tea.Cmd {
	return c.Wrapper(func() (any, error) {
		output, err := c.Client.ListBuckets(ctx, input)
		mssg := c.NewMessage()
		mssg.APIMessage = internal.APIMessage{
			Response: output,
			Err:      err,
		}

		var names []string
		if err == nil {
			for _, b := range output.Buckets {
				names = append(names, *b.Name)
			}

		}
		mssg.Buckets = names
		mssg.LoadBuckets = false
		return mssg, err
	})
}

// func (c *S3Client) PutObject(ctx context.Context, input *s3.PutObjectInput) tea.Cmd {
// 	return c.Wrapper(func() (any, error) {
// 		return c.client.PutObject(ctx, input)
// 	})
// }

// func (c *S3Client) GetObject(ctx context.Context, input *s3.GetObjectInput) tea.Cmd {
// 	return c.Wrapper(func() (any, error) {
// 		return c.client.GetObject(ctx, input)
// 	})
// }

// func (c *S3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) tea.Cmd {
// 	return c.Wrapper(func() (any, error) {
// 		return c.client.DeleteObject(ctx, input)
// 	})
// }
