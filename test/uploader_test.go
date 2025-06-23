package test

import (
	"errors"
	"testing"

	"github.com/anshul543/data-ingestion/internal/config"
	"github.com/anshul543/data-ingestion/internal/fetcher"
	"github.com/anshul543/data-ingestion/internal/transformer"
	"github.com/anshul543/data-ingestion/internal/uploader"
)

func sampleData() []transformer.TransformedPost {
	return []transformer.TransformedPost{
		{
			Post: fetcher.Post{
				UserID: 1,
				ID:     1,
				Title:  "Test",
				Body:   "Hello",
			},
			IngestedAt: "2025-01-01T00:00:00Z",
			Source:     "placeholder_api",
		},
	}
}

// ----- Real S3 Uploader with Invalid Credentials -----

func TestUploader_InvalidAWSConfig_ShouldFail(t *testing.T) {
	cfg := config.AWSConfig{
		AccessKey: "invalid",
		SecretKey: "invalid",
		Region:    "us-east-1",
		Bucket:    "fake-bucket",
	}

	u := uploader.NewUploader(cfg)
	err := u.Upload(sampleData())

	if err == nil {
		t.Error("Expected failure due to invalid AWS credentials, got nil")
	}
}

// ----- Mock Success -----

type mockUploader struct {
	Called bool
	Data   []transformer.TransformedPost
}

func (m *mockUploader) Upload(posts []transformer.TransformedPost) error {
	m.Called = true
	m.Data = posts
	return nil
}

func TestUploader_MockSuccess(t *testing.T) {
	mock := &mockUploader{}
	err := mock.Upload(sampleData())

	if err != nil {
		t.Errorf("Expected success from mock uploader, got: %v", err)
	}
	if !mock.Called {
		t.Error("Expected mock uploader to be called")
	}
	if len(mock.Data) != 1 {
		t.Errorf("Expected 1 post, got %d", len(mock.Data))
	}
}

// ----- Mock Failure -----

type failingMockUploader struct{}

func (f *failingMockUploader) Upload([]transformer.TransformedPost) error {
	return errors.New("mocked upload failure")
}

func TestUploader_MockFailure(t *testing.T) {
	mock := &failingMockUploader{}
	err := mock.Upload(sampleData())

	if err == nil {
		t.Error("Expected error from mock uploader, got nil")
	}
}
