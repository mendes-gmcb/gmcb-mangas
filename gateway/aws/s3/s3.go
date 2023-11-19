package s3

import (
	"trabalho/constants/str"
	"trabalho/env"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Uploader interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

func NewUploader() (Uploader, error) {
	sess, err := createAWSSession()
	if err != nil {
		return nil, err
	}

	uploader := s3manager.NewUploader(sess)
	return uploader, nil
}

func createAWSSession() (*session.Session, error) {
	creds := credentials.NewStaticCredentials(env.AWS_CONFIG.AccessKey, env.AWS_CONFIG.SecretKey, str.EMPTY_STRING)
	config := aws.NewConfig().WithCredentials(creds).WithRegion(env.AWS_CONFIG.Region)

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
