package imbeddings

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestClientBasic(t *testing.T) {
	baseURL := os.Getenv("IMBEDDINGS_BASE_URL")
	if baseURL == "" {
		t.Skip("IMBEDDINGS_BASE_URL not set")
	}
	model := os.Getenv("IMBEDDINGS_MODEL_ID")
	if model == "" {
		model = "facebook/dinov2-small"
	}
	apiKey := os.Getenv("IMBEDDINGS_API_KEY")

	client, err := NewClient(baseURL, model, nil)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if apiKey != "" {
		client.SetAPIKey(apiKey)
	}

	imgPath := filepath.Join("..", "imbeddings", "tests", "images", "img.webp")
	file, err := os.Open(imgPath)
	if err != nil {
		t.Fatalf("open image: %v", err)
	}
	defer file.Close()

	params := ImbeddingsParams{
		Images: []Image{ImageReader(file)},
	}
	vectors, err := client.New(context.Background(), params)
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}
	if len(vectors) == 0 {
		t.Fatalf("expected embeddings")
	}
	if len(vectors[0].Vector) == 0 {
		t.Fatalf("expected embedding vector")
	}
}

func TestClientBase64Output(t *testing.T) {
	baseURL := os.Getenv("IMBEDDINGS_BASE_URL")
	if baseURL == "" {
		t.Skip("IMBEDDINGS_BASE_URL not set")
	}
	model := os.Getenv("IMBEDDINGS_MODEL_ID")
	if model == "" {
		model = "facebook/dinov2-small"
	}
	apiKey := os.Getenv("IMBEDDINGS_API_KEY")

	client, err := NewClient(baseURL, model, nil)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if apiKey != "" {
		client.SetAPIKey(apiKey)
	}

	imgPath := filepath.Join("..", "imbeddings", "tests", "images", "img.webp")
	file, err := os.ReadFile(imgPath)
	if err != nil {
		t.Fatalf("read image: %v", err)
	}
	b64 := base64.StdEncoding.EncodeToString(file)

	params := ImbeddingsParams{
		Images:         []Image{ImageBase64(b64)},
		EncodingFormat: "base64",
	}
	vectors, err := client.New(context.Background(), params)
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}
	if len(vectors) == 0 {
		t.Fatalf("expected embeddings")
	}
	if vectors[0].Base64 == "" {
		t.Fatalf("expected base64 embedding")
	}
}
