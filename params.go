package imbeddings

import (
	"encoding/json"
	"io"
)

// embeddingRequest is the top-level request payload for the embeddings API.
type embeddingRequest struct {
	Input          []string `json:"input"`
	Model          string   `json:"model"`
	EncodingFormat string   `json:"encoding_format,omitempty"`
	Dimensions     *int     `json:"dimensions,omitempty"`
	User           string   `json:"user,omitempty"`
}

// embeddingItem represents a single embeddings result with its index.
type embeddingItem struct {
	Index     int             `json:"index"`
	Embedding json.RawMessage `json:"embedding"`
}

// embeddingUsage describes token usage info from the service.
type embeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// embeddingResponse is the response payload from the embeddings API.
type embeddingResponse struct {
	Data  []embeddingItem `json:"data"`
	Model string          `json:"model"`
	Usage embeddingUsage  `json:"usage"`
}

// Embedding contains a single returned vector or base64 payload.
type Embedding struct {
	Vector []float32
	Base64 string
}

// Image describes a single input image, provided as a URL, base64 string, or reader.
type Image struct {
	URL    string
	Base64 string
	Reader io.Reader
}

// ImbeddingsParams configures an embeddings request.
// Provide either Images or Inputs, not both.
type ImbeddingsParams struct {
	Images         []Image
	Inputs         []string
	Model          string
	EncodingFormat string
	Dimensions     int
	User           string
}

// ImageURL returns an Image configured with a URL source.
func ImageURL(url string) Image {
	return Image{URL: url}
}

// ImageBase64 returns an Image configured with a base64 or data URI source.
func ImageBase64(data string) Image {
	return Image{Base64: data}
}

// ImageReader returns an Image configured with a reader source.
func ImageReader(r io.Reader) Image {
	return Image{Reader: r}
}
