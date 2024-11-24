package r2

import (
	"r2-api-go/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func initS3Session() (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Configs.R2_Region),
		Endpoint:    aws.String(config.Configs.R2_Endpoint),
		Credentials: credentials.NewStaticCredentials(config.Configs.R2_Access_Key, config.Configs.R2_Secret_Key, ""),
	})
	if err != nil {
		return nil, err
	}
	return s3.New(sess), nil
}

func GenerateDownloadURL(key string) (*string, error) {
	svc, err := initS3Session()
	if err != nil {
		return nil, err
	}

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(config.Configs.R2_Bucket),
		Key:    aws.String(key),
	})

	downloadURL, err := req.Presign(config.Configs.R2_Download_Expiry)
	if err != nil {
		return nil, err
	}
	return &downloadURL, nil
}

func GenerateUploadURL(key string) (*string, error) {
	svc, err := initS3Session()
	if err != nil {
		return nil, err
	}

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(config.Configs.R2_Bucket),
		Key:    aws.String(key),
	})

	uploadURL, err := req.Presign(config.Configs.R2_Upload_Expiry)
	if err != nil {
		return nil, err
	}

	return &uploadURL, nil
}
