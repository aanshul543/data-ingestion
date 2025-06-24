package transformer

import (
	"testing"
	"time"

	"github.com/anshul543/data-ingestion/internal/fetcher"
)

func TestTransformPosts_Success(t *testing.T) {
	sample := []fetcher.Post{
		{UserID: 1, ID: 1, Title: "Test", Body: "Hello"},
	}

	transformed := TransformPosts(sample)
	if len(transformed) != 1 {
		t.Errorf("Expected 1 item, got %d", len(transformed))
	}

	if transformed[0].Source != "placeholder_api" {
		t.Errorf("Expected source 'placeholder_api', got '%s'", transformed[0].Source)
	}

	_, err := time.Parse(time.RFC3339, transformed[0].IngestedAt)
	if err != nil {
		t.Errorf("Invalid timestamp format: %v", transformed[0].IngestedAt)
	}
}

func TestTransformPosts_EmptyInput(t *testing.T) {
	var sample []fetcher.Post
	result := TransformPosts(sample)

	if len(result) != 0 {
		t.Errorf("Expected 0 items for empty input, got %d", len(result))
	}
}
