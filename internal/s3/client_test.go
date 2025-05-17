package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

// Mock S3 client
type mockS3 struct {
	ListBucketsFunc   func(ctx context.Context, input *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	CreateBucketFunc  func(ctx context.Context, input *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	ListObjectsV2Func func(ctx context.Context, input *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	PutObjectFunc     func(ctx context.Context, input *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObjectFunc     func(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	DeleteObjectFunc  func(ctx context.Context, input *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	HeadObjectFunc    func(ctx context.Context, input *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

func (m *mockS3) ListBuckets(ctx context.Context, input *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.ListBucketsFunc(ctx, input, optFns...)
}
func (m *mockS3) CreateBucket(ctx context.Context, input *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return m.CreateBucketFunc(ctx, input, optFns...)
}
func (m *mockS3) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	return m.ListObjectsV2Func(ctx, input, optFns...)
}
func (m *mockS3) PutObject(ctx context.Context, input *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m.PutObjectFunc(ctx, input, optFns...)
}
func (m *mockS3) GetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.GetObjectFunc(ctx, input, optFns...)
}
func (m *mockS3) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return m.DeleteObjectFunc(ctx, input, optFns...)
}
func (m *mockS3) HeadObject(ctx context.Context, input *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return m.HeadObjectFunc(ctx, input, optFns...)
}

func TestListBuckets(t *testing.T) {
	mock := &mockS3{
		ListBucketsFunc: func(ctx context.Context, input *s3.ListBucketsInput, _ ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
			return &s3.ListBucketsOutput{
				Buckets: []types.Bucket{{Name: aws.String("bucket1")}},
			}, nil
		},
	}
	client := &S3Client{Client: mock}
	cmd := client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	msg := cmd().(S3MenuMessage)
	assert.Equal(t, S3OpListBuckets, msg.Op)
	assert.Len(t, msg.Buckets, 1)
	assert.Equal(t, "bucket1", *msg.Buckets[0].Name)
}

func TestCreateBucket(t *testing.T) {
	mock := &mockS3{
		CreateBucketFunc: func(ctx context.Context, input *s3.CreateBucketInput, _ ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
			return &s3.CreateBucketOutput{}, nil
		},
	}
	client := &S3Client{Client: mock}
	cmd := client.CreateBucket(context.Background(), &s3.CreateBucketInput{Bucket: aws.String("test-bucket")})
	msg := cmd().(S3MenuMessage)
	assert.Equal(t, S3OpCreateBucket, msg.Op)
	assert.Equal(t, "test-bucket", msg.Bucket)
}

func TestListObjects(t *testing.T) {
	mock := &mockS3{
		ListObjectsV2Func: func(ctx context.Context, input *s3.ListObjectsV2Input, _ ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			return &s3.ListObjectsV2Output{
				Contents: []types.Object{
					{Key: aws.String("file1.txt")},
					{Key: aws.String("file2.txt")},
				},
			}, nil
		},
	}
	client := &S3Client{Client: mock}
	cmd := client.ListObjects(context.Background(), &s3.ListObjectsV2Input{})
	msg := cmd().(S3MenuMessage)
	assert.Equal(t, S3OpListObjects, msg.Op)
	assert.ElementsMatch(t, []string{"file1.txt", "file2.txt"}, msg.Objects)
}

func TestPutObject_FileNotFound(t *testing.T) {
	mock := &mockS3{
		PutObjectFunc: func(ctx context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			return &s3.PutObjectOutput{}, nil
		},
	}
	client := &S3Client{Client: mock}
	cmd := client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	}, "nonexistent.txt")
	msg := cmd().(S3MenuMessage)
	assert.NotNil(t, msg.APIMessage.Err)
}

func TestPutObject_Success(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("hello")
	tmpfile.Close()

	mock := &mockS3{
		PutObjectFunc: func(ctx context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			return &s3.PutObjectOutput{}, nil
		},
	}
	client := &S3Client{Client: mock}
	cmd := client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	}, tmpfile.Name())
	msg := cmd().(S3MenuMessage)
	assert.Equal(t, S3OpPutObject, msg.Op)
	assert.Contains(t, msg.APIMessage.Status, "Uploaded")
}

func TestGetObject_FileWrite(t *testing.T) {
	tmpdir := t.TempDir()
	mock := &mockS3{
		GetObjectFunc: func(ctx context.Context, input *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body:          io.NopCloser(bytes.NewReader([]byte("data"))),
				ContentLength: 4,
			}, nil
		},
	}
	client := &S3Client{Client: mock}
	cmd := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("foo/bar.txt"),
	}, tmpdir)
	msg := cmd().(S3MenuMessage)
	assert.Equal(t, S3OpGetObject, msg.Op)
	assert.Contains(t, msg.APIMessage.Status, "Fetched")
	// Check file exists
	_, err := os.Stat(filepath.Join(tmpdir, "bar.txt"))
	assert.NoError(t, err)
}

func TestGetObject_Error(t *testing.T) {
	mock := &mockS3{
		GetObjectFunc: func(ctx context.Context, input *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return nil, errors.New("fail")
		},
	}
	client := &S3Client{Client: mock}
	cmd := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("foo/bar.txt"),
	}, t.TempDir())
	msg := cmd().(S3MenuMessage)
	assert.NotNil(t, msg.APIMessage.Err)
}

func TestDeleteObject(t *testing.T) {
	mock := &mockS3{
		DeleteObjectFunc: func(ctx context.Context, input *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
			return &s3.DeleteObjectOutput{}, nil
		},
	}
	client := &S3Client{Client: mock}
	cmd := client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	msg := cmd().(S3MenuMessage)
	assert.Equal(t, S3OpDeleteObject, msg.Op)
	assert.Contains(t, msg.APIMessage.Status, "Deleted")
}

func TestGetObjectMetadata(t *testing.T) {
	now := time.Now()
	mock := &mockS3{
		HeadObjectFunc: func(ctx context.Context, input *s3.HeadObjectInput, _ ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
			return &s3.HeadObjectOutput{
				ContentType:   aws.String("text/plain"),
				ContentLength: aws.Int64(123),
				LastModified:  &now,
				ETag:          aws.String("etag"),
				StorageClass:  types.StorageClassStandard,
				Metadata:      map[string]string{"foo": "bar"},
			}, nil
		},
	}
	client := &S3Client{Client: mock}
	cmd := client.GetObjectMetadata(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	})
	msg := cmd().(S3MenuMessage)
	assert.Equal(t, S3OpGetObjectMetadata, msg.Op)
	assert.Equal(t, "key", msg.Metadata.Key)
	assert.Equal(t, "bucket", msg.Metadata.Bucket)
	assert.Equal(t, "text/plain", msg.Metadata.ContentType)
	assert.Equal(t, int64(123), msg.Metadata.ContentLength)
	assert.Equal(t, now, msg.Metadata.LastModified)
	assert.Equal(t, "etag", msg.Metadata.ETag)
	assert.Equal(t, types.StorageClassStandard, msg.Metadata.StorageClass)
	assert.Equal(t, "bar", msg.Metadata.Metadata["foo"])
}

func TestSplitLast(t *testing.T) {
	a, b := splitLast("foo/bar/baz.txt", "/")
	assert.Equal(t, "foo/bar", a)
	assert.Equal(t, "baz.txt", b)

	a, b = splitLast("filename.txt", "/")
	assert.Equal(t, "filename.txt", a)
	assert.Equal(t, "", b)
}
