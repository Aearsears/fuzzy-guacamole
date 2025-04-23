package internal

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	tea "github.com/charmbracelet/bubbletea"
)

type Client struct {
	S3 *s3.Client
}

func (c *Client) PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return c.S3.PutObject(ctx, input)
}

func (c *Client) GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return c.S3.GetObject(ctx, input)
}

func (c *Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return c.S3.DeleteObject(ctx, input)
}
