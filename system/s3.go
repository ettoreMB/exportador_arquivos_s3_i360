package system

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Upload(filepath string, config *Config) error {

	var bucketName = config.S3.Bucket
	var accessKeyId = config.S3.key
	var accessKeySecret = config.S3.Secret

	cfg, err := s3config.LoadDefaultConfig(context.TODO(),
		s3config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		s3config.WithRegion("eu-central-1"),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	parsedName := strings.Split(file.Name(), ".")
	fileName := strings.Split(parsedName[1], "/")
	fullKey := fmt.Sprintf("%s/%s", config.S3.Path, fileName[2])
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fullKey),
		Body:          file,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String("text/plain"),
	})

	if err != nil {
		return err
	}
	return nil
}
