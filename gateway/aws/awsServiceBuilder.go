package aws

import (
	"trabalho/env"

	aws_manager "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSServiceBuilder interface {
	BuildAWSClient() *s3.S3
	BuildUploader() *s3manager.Uploader
	BuildSession() *session.Session
	WithAccessKey(accessKey string) *awsBuilder
	WithSecretKey(secretKey string) *awsBuilder
	WithRegion(region string) *awsBuilder
}

type awsBuilder struct {
	accessKey string
	secretKey string
	region    string
}

func NewDefaultServiceBuilder() AWSServiceBuilder {
	return &awsBuilder{
		accessKey: env.AWS_CONFIG.AccessKey,
		secretKey: env.AWS_CONFIG.SecretKey,
		region: env.AWS_CONFIG.Region,
	}
}

func (builder *awsBuilder) BuildAWSClient() *s3.S3 {
	sess := builder.BuildSession()
	return s3.New(sess)
}

func (builder *awsBuilder) BuildUploader() *s3manager.Uploader {
	sess := builder.BuildSession()
	return s3manager.NewUploader(sess)
}

func (builder *awsBuilder) BuildSession() *session.Session {
	creds := credentials.NewStaticCredentials(builder.accessKey, builder.secretKey, "")
	sess, err := session.NewSession(&aws_manager.Config{
		Credentials: creds,
		Region:      aws_manager.String(builder.region),
	})
	if err != nil {
		panic(err)
	}
	return sess
}

func (builder *awsBuilder) WithAccessKey(accessKey string) *awsBuilder {
	builder.accessKey = accessKey
	return builder
}

func (builder *awsBuilder) WithSecretKey(secretKey string) *awsBuilder {
	builder.secretKey = secretKey
	return builder
}

func (builder *awsBuilder) WithRegion(region string) *awsBuilder {
	builder.region = region
	return builder
}

