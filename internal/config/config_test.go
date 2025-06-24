package config_test

import (
	"os"
	"testing"

	"github.com/anshul543/data-ingestion/internal/config"
)

func TestLoadAWSConfig(t *testing.T) {
	// Setup test environment variables
	os.Setenv("AWS_ACCESS_KEY_ID", "test-access-key")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test-secret-key")
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_BUCKET", "test-bucket")
	os.Setenv("S3_ENDPOINT", "http://localhost:9000")

	// Call the function
	cfg := config.LoadAWSConfig()

	// Validate results
	if cfg.AccessKey != "test-access-key" {
		t.Errorf("Expected AccessKey to be 'test-access-key', got '%s'", cfg.AccessKey)
	}
	if cfg.SecretKey != "test-secret-key" {
		t.Errorf("Expected SecretKey to be 'test-secret-key', got '%s'", cfg.SecretKey)
	}
	if cfg.Region != "us-east-1" {
		t.Errorf("Expected Region to be 'us-east-1', got '%s'", cfg.Region)
	}
	if cfg.Bucket != "test-bucket" {
		t.Errorf("Expected Bucket to be 'test-bucket', got '%s'", cfg.Bucket)
	}
	if cfg.Endpoint != "http://localhost:9000" {
		t.Errorf("Expected Endpoint to be 'http://localhost:9000', got '%s'", cfg.Endpoint)
	}
}
