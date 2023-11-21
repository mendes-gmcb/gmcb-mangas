package aws

import "github.com/aws/aws-sdk-go/service/s3"

func NewSession() *s3.S3 {
	serviceBuilder := NewDefaultServiceBuilder()
	return serviceBuilder.BuildAWSClient()
}