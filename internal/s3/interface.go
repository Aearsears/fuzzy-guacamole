package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	tea "github.com/charmbracelet/bubbletea"
)

type S3API interface {
	PutObject(ctx context.Context, input *s3.PutObjectInput, filePath string) tea.Cmd
	GetObject(ctx context.Context, input *s3.GetObjectInput, savePath string) tea.Cmd
	GetObjectMetadata(ctx context.Context, input *s3.HeadObjectInput) tea.Cmd
	DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) tea.Cmd
	ListBuckets(ctx context.Context, input *s3.ListBucketsInput) tea.Cmd
	CreateBucket(ctx context.Context, input *s3.CreateBucketInput) tea.Cmd
	ListObjects(ctx context.Context, input *s3.ListObjectsV2Input) tea.Cmd
}
