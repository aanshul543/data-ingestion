package uploader

import "github.com/anshul543/data-ingestion/internal/transformer"

type S3Uploader interface {
	Upload(transformed []transformer.TransformedPost) error
}
