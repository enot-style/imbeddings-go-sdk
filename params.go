package imbeddings

import "io"

// embeddingInputItem is a single image input for the API payload.
type embeddingInputItem struct {
	Type        string `json:"type"`
	ImageBase64 string `json:"image_base64,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

// embeddingRequest is the top-level request payload for the embeddings API.
type embeddingRequest struct {
	Input []embeddingInputItem `json:"input"`
	Model string               `json:"model"`
}

// embeddingVectors holds the mean and cls vectors returned by the service.
type embeddingVectors struct {
	Cls  []float32 `json:"cls"`
	Mean []float32 `json:"mean"`
}

// embeddingItem represents a single embeddings result with its index.
type embeddingItem struct {
	Index      int              `json:"index"`
	Embeddings embeddingVectors `json:"embeddings"`
}

// embeddingResponse is the response payload from the embeddings API.
type embeddingResponse struct {
	Data []embeddingItem `json:"data"`
}

// Embeddings contains the returned mean and cls vectors.
type Embeddings struct {
	Mean []float32
	Cls  []float32
}

// Image describes a single input image, provided as a URL or reader.
type Image struct {
	URL    string
	Reader io.Reader
}

// ImbeddingsParams configures a embeddings request.
type ImbeddingsParams struct {
	Images []Image
	Model  string
}

// ImageURL returns an Image configured with a URL source.
func ImageURL(url string) Image {
	return Image{URL: url}
}

// ImageReader returns an Image configured with a reader source.
func ImageReader(r io.Reader) Image {
	return Image{Reader: r}
}
