package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Aearsears/fuzzy-guacamole/internal"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	tea "github.com/charmbracelet/bubbletea"
)

type S3Client struct {
	Client *s3.Client
}
type S3OperationType int

const (
	S3OpListBuckets S3OperationType = iota
	S3OpCreateBucket
	S3OpListObjects
	S3OpGetObject
	S3OpPutObject
	S3OpDeleteObject
	S3OpGetObjectMetadata
)

type S3ObjectMetadata struct {
	Key           string
	Bucket        string
	ContentType   string
	ContentLength int64
	LastModified  time.Time
	ETag          string
	StorageClass  types.StorageClass
	Metadata      map[string]string
}

type S3MenuMessage struct {
	Op         S3OperationType
	APIMessage internal.APIMessage
	Buckets    []types.Bucket // for ListBuckets
	Objects    []string       // for ListObjects
	Bucket     string
	Metadata   S3ObjectMetadata
}

func (c *S3Client) NewMessage() S3MenuMessage {
	return S3MenuMessage{
		APIMessage: internal.APIMessage{},
		Buckets:    nil,
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
		mssg.Op = S3OpListBuckets

		if err != nil {
			return mssg, err
		}

		mssg.Buckets = output.Buckets
		return mssg, err
	})
}

func (c *S3Client) CreateBucket(ctx context.Context, input *s3.CreateBucketInput) tea.Cmd {
	return c.Wrapper(func() (any, error) {
		output, err := c.Client.CreateBucket(ctx, input)
		mssg := c.NewMessage()
		mssg.APIMessage = internal.APIMessage{
			Response: output,
			Err:      err,
		}
		mssg.Op = S3OpCreateBucket
		mssg.Bucket = *input.Bucket

		if err != nil {
			return mssg, err
		}

		return mssg, err
	})
}

func (c *S3Client) ListObjects(ctx context.Context, input *s3.ListObjectsV2Input) tea.Cmd {
	return c.Wrapper(func() (any, error) {
		resp, err := c.Client.ListObjectsV2(ctx, input)
		mssg := c.NewMessage()
		mssg.APIMessage = internal.APIMessage{
			Response: resp,
			Err:      err,
		}
		mssg.Op = S3OpListObjects

		if err != nil {
			return mssg, err
		}

		var objs []string
		for _, obj := range resp.Contents {
			objs = append(objs, *obj.Key)
		}
		mssg.Objects = objs

		return mssg, err
	})
}

// PutObject uploads a file to S3. filePath is relative to the current working directory of the TUI
func (c *S3Client) PutObject(ctx context.Context, input *s3.PutObjectInput, filePath string) tea.Cmd {
	return c.Wrapper(func() (any, error) {
		//TODO: handle large objects

		mssg := c.NewMessage()
		mssg.APIMessage = internal.APIMessage{}
		file, err := os.Open(filePath)
		if err != nil {
			mssg.APIMessage.Err = err
			return mssg, err
		}

		defer file.Close()
		input.Body = file

		resp, err := c.Client.PutObject(ctx, input)
		mssg.APIMessage.Response = resp
		mssg.APIMessage.Err = err
		mssg.APIMessage.Status = fmt.Sprintf("Uploaded %s/%s successfully", *input.Bucket, *input.Key)
		mssg.Op = S3OpPutObject

		return mssg, err
	})
}

func (c *S3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, savePath string) tea.Cmd {
	return c.Wrapper(func() (any, error) {
		//TODO: handle large objects
		resp, err := c.Client.GetObject(ctx, input)
		mssg := c.NewMessage()
		mssg.APIMessage = internal.APIMessage{
			Response: resp,
			Err:      err,
			Status:   fmt.Sprintf("Fetched %s/%s successfully", *input.Bucket, *input.Key),
		}
		mssg.Op = S3OpGetObject

		if err != nil {
			return mssg, err
		}

		defer resp.Body.Close()

		_, tail := splitLast(*input.Key, "/")
		outFile, err := os.Create(filepath.Join(savePath, tail))
		if err != nil {
			mssg.APIMessage.Err = err
			return mssg, err
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			mssg.APIMessage.Err = err
			return mssg, err
		}

		return mssg, err
	})
}

// splitLast splits a string s into two parts at the last occurrence of sep
func splitLast(s, sep string) (string, string) {
	idx := strings.LastIndex(s, sep)
	if idx == -1 {
		return s, "" // sep not found
	}
	return s[:idx], s[idx+len(sep):]
}

func (c *S3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput) tea.Cmd {
	return c.Wrapper(func() (any, error) {
		resp, err := c.Client.DeleteObject(ctx, input)
		mssg := c.NewMessage()
		mssg.APIMessage = internal.APIMessage{
			Response: resp,
			Err:      err,
		}
		mssg.Op = S3OpDeleteObject
		if err != nil {
			return mssg, err
		}

		mssg.APIMessage.Status = fmt.Sprintf("Deleted %s/%s successfully", *input.Bucket, *input.Key)
		return mssg, err
	})
}

func (c *S3Client) GetObjectMetadata(ctx context.Context, input *s3.HeadObjectInput) tea.Cmd {
	return c.Wrapper(func() (any, error) {
		resp, err := c.Client.HeadObject(ctx, input)
		mssg := c.NewMessage()
		mssg.APIMessage = internal.APIMessage{
			Response: resp,
			Err:      err,
		}
		mssg.Op = S3OpGetObjectMetadata

		if err != nil {
			return mssg, err
		}

		metadata := S3ObjectMetadata{
			Key:           *input.Key,
			Bucket:        *input.Bucket,
			ContentType:   aws.ToString(resp.ContentType),
			ContentLength: *resp.ContentLength,
			LastModified:  *resp.LastModified,
			ETag:          aws.ToString(resp.ETag),
			StorageClass:  resp.StorageClass,
			Metadata:      resp.Metadata,
		}
		mssg.Metadata = metadata

		mssg.APIMessage.Status = fmt.Sprintf("Fetched %s/%s metadata successfully", *input.Bucket, *input.Key)
		return mssg, err

	})

}
