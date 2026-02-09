package imbeddings

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ImbeddingsClient provides a minimal SDK for the imbeddings service.
type ImbeddingsClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewClient initializes a client with the given base URL and model.
// If httpClient is nil, DefaultClient is used.
func NewClient(baseURL, model string, httpClient *http.Client) (*ImbeddingsClient, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, ValidationError{Op: ValidationOpInit, Err: fmt.Errorf("missing base URL")}
	}
	if strings.TrimSpace(model) == "" {
		return nil, ValidationError{Op: ValidationOpInit, Err: fmt.Errorf("missing model name")}
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &ImbeddingsClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		model:      model,
		httpClient: httpClient,
	}, nil
}

func (c *ImbeddingsClient) Model() string {
	return c.model
}

func (c *ImbeddingsClient) SetModel(model string) error {
	if strings.TrimSpace(model) == "" {
		return ValidationError{Op: ValidationOpSetModel, Err: fmt.Errorf("missing model name")}
	}
	c.model = model
	return nil
}

func (c *ImbeddingsClient) New(ctx context.Context, params ImbeddingsParams) ([]Embeddings, error) {
	if len(params.Images) == 0 {
		return nil, ValidationError{Op: ValidationOpParams, Err: fmt.Errorf("missing images")}
	}
	if ctx == nil {
		ctx = context.Background()
	}

	inputs := make([]embeddingInputItem, 0, len(params.Images))
	for i, img := range params.Images {
		item, err := imageToInputItem(img)
		if err != nil {
			return nil, ValidationError{Op: ValidationOpImage, Err: fmt.Errorf("index %d: %w", i, err)}
		}
		inputs = append(inputs, item)
	}

	model := strings.TrimSpace(params.Model)
	if model == "" {
		model = c.model
	}

	payload := embeddingRequest{
		Input: inputs,
		Model: model,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, RequestError{Op: RequestOpMarshalRequest, Err: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, RequestError{Op: RequestOpCreateRequest, Err: err}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, RequestError{Op: RequestOpCallService, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, RequestError{
			Op:  "http status",
			Err: fmt.Errorf("%s: %s", resp.Status, strings.TrimSpace(string(respBody))),
		}
	}

	var parsed embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, ResponseError{Op: ResponseOpDecodeResponse, Err: err}
	}
	if len(parsed.Data) == 0 {
		return nil, ResponseError{Op: ResponseOpEmptyResponse, Err: fmt.Errorf("no embeddings")}
	}

	results := make([]Embeddings, 0, len(parsed.Data))
	for i, item := range parsed.Data {
		result := Embeddings{
			Mean: item.Embeddings.Mean,
			Cls:  item.Embeddings.Cls,
		}
		if len(result.Mean) == 0 && len(result.Cls) == 0 {
			return nil, ResponseError{
				Op:  ResponseOpResult,
				Err: fmt.Errorf("index %d: no vectors", i),
			}
		}
		results = append(results, result)
	}
	return results, nil
}

func imageToInputItem(img Image) (embeddingInputItem, error) {
	if strings.TrimSpace(img.URL) != "" {
		return embeddingInputItem{
			Type:     "image",
			ImageURL: img.URL,
		}, nil
	}
	if img.Reader == nil {
		return embeddingInputItem{}, ValidationError{Op: ValidationOpImage, Err: fmt.Errorf("missing image reader")}
	}

	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	if _, err := io.Copy(enc, img.Reader); err != nil {
		return embeddingInputItem{}, RequestError{Op: RequestOpReadImageData, Err: err}
	}
	if err := enc.Close(); err != nil {
		return embeddingInputItem{}, RequestError{Op: RequestOpEncodeImage, Err: err}
	}

	return embeddingInputItem{
		Type:        "image",
		ImageBase64: buf.String(),
	}, nil
}
