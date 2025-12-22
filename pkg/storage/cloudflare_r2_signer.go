package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jorgeAM/go-template/pkg/errors"
)

var _ Signer = (*CloudflareR2Signer)(nil)

type CloudflareR2Signer struct {
	s3Client   *s3.Client
	bucketName string
}

func NewCloudflareR2Client(bucketName, accessKey, secretKey, endpoint string) (*CloudflareR2Signer, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &CloudflareR2Signer{
		s3Client:   s3Client,
		bucketName: bucketName,
	}, nil
}

func (c *CloudflareR2Signer) GeneratePresignedURL(ctx context.Context, filename string, contentType ContentType) (string, error) {
	presignClient := s3.NewPresignClient(c.s3Client)
	params := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(filename),
		ContentType: aws.String(contentType.String()),
	}

	presignedURL, err := presignClient.PresignPutObject(ctx, params, s3.WithPresignExpires(5*time.Minute))
	if err != nil {
		return "", errors.Wrap(ErrStorageInternal, err, "there was an error try to generate pre signed url")
	}

	return presignedURL.URL, nil
}
