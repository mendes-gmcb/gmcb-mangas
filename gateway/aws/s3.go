package aws

import "github.com/aws/aws-sdk-go/service/s3/s3manager"

type Uploader interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

func NewUploader() (Uploader) {
	serviceBuilder := NewDefaultServiceBuilder()
	return serviceBuilder.BuildUploader()
}
