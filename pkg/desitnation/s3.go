package desitnation

import (
	"bytes"
	"context"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

type S3 struct {
	cl     *s3.S3
	bucket string
}

func NewS3() (*S3, error) {
	config := &aws.Config{
		S3ForcePathStyle: aws.Bool(true),
	}

	if v := os.Getenv("S3_ACCESS_KEY"); v != "" {
		config.Credentials = credentials.NewStaticCredentials(v, os.Getenv("S3_SECRET_KEY"), "")
	}
	if v := os.Getenv("S3_ENDPOINT"); v != "" {
		config.Endpoint = aws.String(v)
	}
	if v := os.Getenv("S3_REGION"); v != "" {
		config.Region = aws.String(v)
	}
	if v, _ := strconv.ParseBool(os.Getenv("S3_DISABLE_SSL")); v {
		config.DisableSSL = aws.Bool(v)
	}

	ses, err := session.NewSession(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &S3{
		cl:     s3.New(ses),
		bucket: os.Getenv("S3_BUCKET"),
	}, nil
}

func (u *S3) Upload(
	ctx context.Context,
	filePath string,
	data []byte,
) error {
	_, err := u.cl.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Key:    aws.String(filePath),
		Bucket: aws.String(u.bucket),
	})

	return err
}
