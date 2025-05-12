package internal

import "github.com/aws/aws-sdk-go-v2/aws"

type APIMessage struct {
	Err      error
	Response any
	Status   string
}

type AWSConfigMessage struct {
	Config aws.Config
}
