package initializers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var Client *s3.S3
var Uploader *s3manager.Uploader

func ConectToS3() (*s3.S3, error) {
	conf := getAWS()

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(conf.AcessKey, conf.SecretKey, ""),
		Region:      aws.String(conf.Region),
	}))
	Client = s3.New(sess)

	return Client, nil
}

func ConectUploaderToS3() (*s3manager.Uploader, error) {
	conf := getAWS()

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(conf.AcessKey, conf.SecretKey, ""),
		Region:      aws.String(conf.Region),
	})
	if err != nil {
		panic(err)
	}

	Uploader = s3manager.NewUploader(sess)

	return Uploader, nil
}
