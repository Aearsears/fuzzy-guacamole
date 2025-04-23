package internal

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type API interface {
	PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	ListBuckets(ctx context.Context, input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
}
