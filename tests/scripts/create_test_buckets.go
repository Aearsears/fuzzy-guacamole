package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomPath() string {
	levels := rand.Intn(3) + 1 // 1 to 3 levels deep
	path := ""
	for range levels {
		path += randomString(rand.Intn(5)+3) + "/"
	}
	path += "file_" + randomString(5) + ".txt"
	return path
}

func randomContent() string {
	return "Random content: " + randomString(rand.Intn(50)+20)
}

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
		o.UsePathStyle = true
	})

	for i := 1; i <= 5; i++ {
		bucketName := fmt.Sprintf("test-bucket-%d", i)

		// Create the bucket
		_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil && !isBucketAlreadyOwned(err) {
			log.Fatalf("failed to create bucket %s: %v", bucketName, err)
		}
		fmt.Println("✔ Bucket ready:", bucketName)

		// Upload 10 random files
		for j := 0; j < 10; j++ {
			key := randomPath()
			content := randomContent()

			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
				Body:   bytes.NewReader([]byte(content)),
			})
			if err != nil {
				log.Printf("❌ Failed to upload to %s/%s: %v", bucketName, key, err)
			} else {
				fmt.Printf("📁 Uploaded to %s: %s\n", bucketName, key)
			}
		}
	}
}

func isBucketAlreadyOwned(err error) bool {
	return err != nil && ( // Avoid failure if bucket exists
	err.Error() == "BucketAlreadyOwnedByYou" ||
		err.Error() == "BucketAlreadyExists")
}
