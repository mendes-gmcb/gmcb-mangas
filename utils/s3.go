package utils

import (
	"os"
	"trabalho/initializers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadCoverImageToS3(coverImagePath string) error {
	file, err := os.Open(coverImagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	s3Object := s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(coverImagePath),
		Body:   file,
	}

	_, err = initializers.Uploader.Upload(&s3Object)
	return err
}

func DeleteCoverImageFromS3(coverImagePath string) error {
	s3Object := s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(coverImagePath),
	}

	_, err := initializers.Client.DeleteObject(&s3Object)
	return err

}
