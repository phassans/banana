package cloudinary

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type (
	client struct {
		logger zerolog.Logger
	}

	Client interface {
		Upload(values map[string]io.Reader) (Response, error)
		MustOpen(f string) (*os.File, error)
	}

	Response struct {
		PublicID          string        `json:"public_id"`
		Version           int           `json:"version"`
		Signature         string        `json:"signature"`
		Width             int           `json:"width"`
		Height            int           `json:"height"`
		Format            string        `json:"format"`
		ResourceType      string        `json:"resource_type"`
		CreatedAt         time.Time     `json:"created_at"`
		Tags              []interface{} `json:"tags"`
		Bytes             int           `json:"bytes"`
		Type              string        `json:"type"`
		Etag              string        `json:"etag"`
		Placeholder       bool          `json:"placeholder"`
		URL               string        `json:"url"`
		SecureURL         string        `json:"secure_url"`
		AccessMode        string        `json:"access_mode"`
		Existing          bool          `json:"existing"`
		OriginalFilename  string        `json:"original_filename"`
		OriginalExtension string        `json:"original_extension"`
	}
)

// NewCloudinaryClient returns a new cloudinary client
func NewCloudinaryClient(logger zerolog.Logger) Client {
	return &client{logger}
}
