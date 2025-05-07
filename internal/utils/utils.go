package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Aearsears/fuzzy-guacamole/internal/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	tea "github.com/charmbracelet/bubbletea"
)

type Client interface {
	Wrapper(func() (any, error)) tea.Cmd
	// NewMessage() T
}

func ClientFactory(clientType string) Client {
	if clientType == "s3" {
		// Endpoint: http://localhost:4566
		// Region: us-east-1
		// Access key: test
		// Secret key: test
		// Custom endpoint resolver
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == awss3.ServiceID {
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

		s3Client := awss3.NewFromConfig(cfg, func(o *awss3.Options) {
			o.BaseEndpoint = aws.String("http://localhost:4566")
			o.UsePathStyle = true
		})
		return &s3.S3Client{Client: s3Client}

	}

	return nil
}

// todo: check types
// func FormatMetadata(meta smithyhttp.ResponseMetadata) string {
// 	return fmt.Sprintf("Request ID: %s, Status: %d", meta.RequestID, meta.HTTPStatusCode)
// }

// wrapper helper function to send messages
func SendMessage(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
func Debug(msg string) {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	f.WriteString(msg + "\n")
}
