package r2

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

type Config struct {
	R2_Endpoint        string        `yaml:"R2_ENDPOINT"`
	R2_Bucket          string        `yaml:"R2_BUCKET"`
	R2_Region          string        `yaml:"R2_REGION"`
	R2_Access_Key      string        `yaml:"R2_ACCESS_KEY"`
	R2_Secret_Key      string        `yaml:"R2_SECRET_KEY"`
	R2_Upload_Expiry   time.Duration `yaml:"R2_UPLOAD_EXPIRY_MINUTES"`
	R2_Download_Expiry time.Duration `yaml:"R2_DOWNLOAD_EXPIRY_MINUTES"`
}

var config Config

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	uplExpiry := time.Duration(viper.GetInt("R2_UPLOAD_EXPIRY_MINUTES")) * time.Minute
	downExpiry := time.Duration(viper.GetInt("R2_DOWNLOAD_EXPIRY_MINUTES")) * time.Minute

	config = Config{
		R2_Endpoint:        viper.GetString("R2_ENDPOINT"),
		R2_Bucket:          viper.GetString("R2_BUCKET"),
		R2_Region:          viper.GetString("R2_REGION"),
		R2_Access_Key:      viper.GetString("R2_ACCESS_KEY"),
		R2_Secret_Key:      viper.GetString("R2_SECRET_KEY"),
		R2_Upload_Expiry:   uplExpiry,
		R2_Download_Expiry: downExpiry,
	}
}

func initS3Session() (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.R2_Region),
		Endpoint:    aws.String(config.R2_Endpoint),
		Credentials: credentials.NewStaticCredentials(config.R2_Access_Key, config.R2_Secret_Key, ""),
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
		Bucket: aws.String(config.R2_Bucket),
		Key:    aws.String(key),
	})

	downloadURL, err := req.Presign(config.R2_Download_Expiry)
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
		Bucket: aws.String(config.R2_Bucket),
		Key:    aws.String(key),
	})

	uploadURL, err := req.Presign(config.R2_Upload_Expiry)
	if err != nil {
		return nil, err
	}

	return &uploadURL, nil
}
