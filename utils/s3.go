package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sync"
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

	err = upload(file, coverImagePath)
	return err
}

func UploadMultipleImagesToS3(
	file *multipart.FileHeader,
	filepath string,
	wg *sync.WaitGroup,
	semaphore <-chan struct{},
) {
	defer wg.Done()

	fileO, err := file.Open()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		<-semaphore
		return
	}

	err = upload(fileO, filepath)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	<-semaphore
}

func upload(file io.Reader, filepath string) error {
	s3Object := s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(filepath),
		Body:   file,
	}

	_, err := initializers.Uploader.Upload(&s3Object)
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

func DeleteFileFromS3(
	coverImagePath string,
	wg *sync.WaitGroup,
	semaphore <-chan struct{},
) {
	defer wg.Done()
	s3Object := s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(coverImagePath),
	}

	_, err := initializers.Client.DeleteObject(&s3Object)
	fmt.Printf("ERROR: %v", err)
	<-semaphore
}
