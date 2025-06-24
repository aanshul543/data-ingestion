package main

import (
	"log"
	"net/http"
	"os"

	"github.com/anshul543/data-ingestion/internal/api"
	"github.com/anshul543/data-ingestion/internal/config"
	"github.com/anshul543/data-ingestion/internal/fetcher"
	"github.com/anshul543/data-ingestion/internal/transformer"
	"github.com/anshul543/data-ingestion/internal/uploader"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	awsCfg := config.LoadAWSConfig()

	if len(os.Args) > 1 && os.Args[1] == "ingest" {
		runIngestion(awsCfg)
		return
	}

	runServer(awsCfg)
}

func runIngestion(awsCfg config.AWSConfig) {
	posts, err := fetcher.FetchPosts()
	if err != nil {
		log.Fatalf("Failed to fetch posts: %v", err)
	}
	transformed := transformer.TransformPosts(posts)
	s3Uploader, err := uploader.NewUploader(awsCfg)
	if err != nil {
		log.Fatalf("Uploader initialization failed: %v", err)
	}
	err = s3Uploader.Upload(transformed)
	if err != nil {
		log.Fatalf("Failed to upload to S3: %v", err)
	}
	log.Println("Ingestion complete.")
}

func runServer(awsCfg config.AWSConfig) {
	httpAPI := api.NewAPI(awsCfg)
	http.HandleFunc("/ingest", httpAPI.IngestData)
	http.HandleFunc("/ingested-posts", httpAPI.GetIngestedPosts)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
