package uploader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/anshul543/data-ingestion/internal/config"
	"github.com/anshul543/data-ingestion/internal/transformer"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// uploader implements S3Uploader interface using real AWS S3
type uploader struct {
	bucket string
	client *s3.Client
}

// NewUploader returns a real S3Uploader implementation or error
func NewUploader(cfg config.AWSConfig) (S3Uploader, error) {
	ctx := context.TODO()

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var s3Client *s3.Client
	if cfg.Endpoint != "" {
		log.Println("⚙️ Using custom S3 endpoint:", cfg.Endpoint)
		s3Client = s3.New(s3.Options{
			Credentials:      awsCfg.Credentials,
			Region:           cfg.Region,
			EndpointResolver: s3.EndpointResolverFromURL(cfg.Endpoint),
			UsePathStyle:     true,
		})
	} else {
		s3Client = s3.NewFromConfig(awsCfg)
	}

	// Attempt to create bucket if not exists
	_, err = s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &cfg.Bucket,
	})
	if err != nil {
		log.Printf("📦 Bucket %s not found. Creating...", cfg.Bucket)
		_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: &cfg.Bucket,
			CreateBucketConfiguration: &s3types.CreateBucketConfiguration{
				LocationConstraint: s3types.BucketLocationConstraint(cfg.Region),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("✅ Bucket created: %s", cfg.Bucket)
	} else {
		log.Printf("✔ Bucket exists: %s", cfg.Bucket)
	}

	return &uploader{
		bucket: cfg.Bucket,
		client: s3Client,
	}, nil
}

// Upload pushes transformed posts to S3 as a JSON file
func (u *uploader) Upload(posts []transformer.TransformedPost) error {
	ctx := context.TODO()

	if u.client == nil {
		return fmt.Errorf("uploader client is nil")
	}

	jsonData, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	key := fmt.Sprintf("ingested/posts_%d.json", time.Now().Unix())

	_, err = u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(jsonData),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	log.Printf("✅ Uploaded file to S3: %s/%s", u.bucket, key)
	return nil
}
