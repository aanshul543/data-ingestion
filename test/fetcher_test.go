package test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/anshul543/data-ingestion/internal/fetcher"
)

// ----- Success Test (Mock HTTP Response) -----

type mockSuccessTransport struct{}

func (m *mockSuccessTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	data := `[{"userId":1,"id":1,"title":"mock title","body":"mock body"}]`
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(data)),
		Header:     make(http.Header),
	}, nil
}

func TestFetchPosts_Success(t *testing.T) {
	client := &http.Client{Transport: &mockSuccessTransport{}}
	posts, err := fetcher.FetchPostsWithClient(client)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(posts) != 1 || posts[0].Title != "mock title" {
		t.Errorf("Unexpected post: %+v", posts)
	}
}

// ----- Failure Test (HTTP client returns error) -----

type mockFailTransport struct{}

func (m *mockFailTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("simulated HTTP failure")
}

func TestFetchPosts_Failure(t *testing.T) {
	client := &http.Client{Transport: &mockFailTransport{}}
	_, err := fetcher.FetchPostsWithClient(client)

	if err == nil {
		t.Error("Expected error from failing transport, got nil")
	}
}
