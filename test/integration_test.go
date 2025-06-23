package test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/anshul543/data-ingestion/internal/fetcher"
	"github.com/anshul543/data-ingestion/internal/transformer"
)

// -------------------------
// Mock HTTP Client
// -------------------------

type mockHTTPTransport struct{}

func (m *mockHTTPTransport) RoundTrip(*http.Request) (*http.Response, error) {
	mockJSON := `[{"userId":1,"id":1,"title":"integration","body":"test post"}]`
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(mockJSON)),
		Header:     make(http.Header),
	}, nil
}

// -------------------------
// Mock S3 Uploader
// -------------------------

type mockS3Uploader struct {
	Called bool
	Posts  []transformer.TransformedPost
}

func (m *mockS3Uploader) Upload(posts []transformer.TransformedPost) error {
	m.Called = true
	m.Posts = posts
	return nil
}

// -------------------------
// Failing Mock S3 Uploader
// -------------------------

type mockFailingS3Uploader struct{}

func (m *mockFailingS3Uploader) Upload([]transformer.TransformedPost) error {
	return errors.New("simulated upload error")
}

// -------------------------
// Test: Success Scenario
// -------------------------

func TestFullPipeline_Success(t *testing.T) {
	client := &http.Client{Transport: &mockHTTPTransport{}}
	posts, err := fetcher.FetchPostsWithClient(client)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	transformed := transformer.TransformPosts(posts)
	if len(transformed) != 1 {
		t.Fatalf("Expected 1 post, got %d", len(transformed))
	}

	mockUploader := &mockS3Uploader{}
	err = mockUploader.Upload(transformed)
	if err != nil {
		t.Fatalf("Mock upload failed: %v", err)
	}

	if !mockUploader.Called {
		t.Error("Expected Upload to be called")
	}
	if mockUploader.Posts[0].Source != "placeholder_api" {
		t.Error("Expected source to be placeholder_api")
	}
	_, err = time.Parse(time.RFC3339, mockUploader.Posts[0].IngestedAt)
	if err != nil {
		t.Errorf("Invalid timestamp: %v", err)
	}
}

// -------------------------
// Test: Upload Failure
// -------------------------

func TestFullPipeline_UploadFailure(t *testing.T) {
	client := &http.Client{Transport: &mockHTTPTransport{}}
	posts, err := fetcher.FetchPostsWithClient(client)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	transformed := transformer.TransformPosts(posts)

	mockUploader := &mockFailingS3Uploader{}
	err = mockUploader.Upload(transformed)
	if err == nil {
		t.Error("Expected upload error, got nil")
	}
}
