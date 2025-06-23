package transformer

import (
	"time"

	"github.com/anshul543/data-ingestion/internal/fetcher"
)

type TransformedPost struct {
	fetcher.Post
	IngestedAt string `json:"ingested_at"`
	Source     string `json:"source"`
}

func TransformPosts(posts []fetcher.Post) []TransformedPost {
	var transformed []TransformedPost
	now := time.Now().UTC().Format(time.RFC3339)

	for _, p := range posts {
		transformed = append(transformed, TransformedPost{
			Post:       p,
			IngestedAt: now,
			Source:     "placeholder_api",
		})
	}
	return transformed
}
