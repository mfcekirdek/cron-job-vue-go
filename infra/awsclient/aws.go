package awsclient

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
	"log"
	"mime/multipart"
	"net/url"
)

type AwsClient interface {
	DeleteFileInS3(fileName string) error
	UploadToS3(fileName, fileType string, file multipart.File) (string, error)
	DetermineS3ImageUrl(fileName string) string
}

type awsClient struct {
	session *session.Session
}

func NewAwsClient() AwsClient {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(viper.GetString("RM_AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(
				viper.GetString("RM_AWS_ACCESS_KEY"),
				viper.GetString("RM_AWS_SECRET_KEY"),
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		log.Fatal("Error when trying to connect aws")
	}

	return awsClient{session: sess}
}

func (client awsClient) DeleteFileInS3(fileName string) error {
	svc := s3.New(client.session)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(viper.GetString("RM_AWS_BUCKET_NAME")),
		Key:    aws.String(fileName),
	}
	_, err := svc.DeleteObject(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (client awsClient) UploadToS3(fileName, fileType string, file multipart.File) (string, error) {
	uploader := s3manager.NewUploader(client.session)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(viper.GetString("RM_AWS_BUCKET_NAME")),
		ACL:         aws.String("public-read"),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(fileType),
	})

	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to upload image on s3 %v", err))
	}

	filePath := client.DetermineS3ImageUrl(fileName)
	return filePath, nil
}

func (client awsClient) DetermineS3ImageUrl(fileName string) string {
	bucketName := viper.GetString("RM_AWS_BUCKET_NAME")
	region := viper.GetString("RM_AWS_REGION")
	filePath := "https://" + bucketName + "." + "s3-" + region + ".amazonaws.com/" + url.QueryEscape(fileName)
	return filePath
}
