package api

import (
	"context"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/anshul543/data-ingestion/internal/config"
	"github.com/anshul543/data-ingestion/internal/fetcher"
	"github.com/anshul543/data-ingestion/internal/transformer"
	"github.com/anshul543/data-ingestion/internal/uploader"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type API struct {
	cfg      config.AWSConfig
	s3Client *s3.Client
	uploader uploader.S3Uploader
}

func NewAPI(cfg config.AWSConfig) *API {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(cfg.Region))
	if err != nil {
		log.Fatalf("Failed to load AWS SDK config: %v", err)
	}

	var s3Client *s3.Client
	if cfg.Endpoint != "" {
		log.Println("⚙️ Using custom S3 endpoint in API:", cfg.Endpoint)

		s3Client = s3.New(s3.Options{
			Credentials:      awsCfg.Credentials,
			Region:           cfg.Region,
			EndpointResolver: s3.EndpointResolverFromURL(cfg.Endpoint),
			UsePathStyle:     true, // ✅ Required for MinIO
		})
	} else {
		s3Client = s3.NewFromConfig(awsCfg)
	}

	return &API{
		cfg:      cfg,
		s3Client: s3Client,
		uploader: uploader.NewUploader(cfg),
	}
}

func (a *API) GetIngestedPosts(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// List objects in the bucket under "ingested/"
	resp, err := a.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(a.cfg.Bucket),
		Prefix: aws.String("ingested/"),
	})
	if err != nil {
		http.Error(w, "Failed to list S3 objects", http.StatusInternalServerError)
		return
	}

	if len(resp.Contents) == 0 {
		http.Error(w, "No ingested data found", http.StatusNotFound)
		return
	}

	// Sort by LastModified to get the latest file
	sort.Slice(resp.Contents, func(i, j int) bool {
		return resp.Contents[i].LastModified.After(*resp.Contents[j].LastModified)
	})
	latest := resp.Contents[0].Key

	// Get the latest object
	obj, err := a.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.cfg.Bucket),
		Key:    latest,
	})
	if err != nil {
		http.Error(w, "Failed to get latest ingested file", http.StatusInternalServerError)
		return
	}
	defer obj.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, obj.Body)
}

func (a *API) IngestData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	posts, err := fetcher.FetchPosts()
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	transformed := transformer.TransformPosts(posts)

	err = a.uploader.Upload(transformed)
	if err != nil {
		http.Error(w, "Failed to upload to S3", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ingestion completed successfully"))
}
