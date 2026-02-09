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
	apiKey     string
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

func (c *ImbeddingsClient) SetAPIKey(apiKey string) {
	c.apiKey = strings.TrimSpace(apiKey)
}

func (c *ImbeddingsClient) New(ctx context.Context, params ImbeddingsParams) ([]Embedding, error) {
	if len(params.Images) == 0 && len(params.Inputs) == 0 {
		return nil, ValidationError{Op: ValidationOpParams, Err: fmt.Errorf("missing inputs")}
	}
	if len(params.Images) > 0 && len(params.Inputs) > 0 {
		return nil, ValidationError{Op: ValidationOpParams, Err: fmt.Errorf("use either Images or Inputs, not both")}
	}
	if ctx == nil {
		ctx = context.Background()
	}

	inputs := make([]string, 0)
	if len(params.Inputs) > 0 {
		for i, input := range params.Inputs {
			if strings.TrimSpace(input) == "" {
				return nil, ValidationError{Op: ValidationOpParams, Err: fmt.Errorf("index %d: empty input", i)}
			}
			inputs = append(inputs, input)
		}
	} else {
		inputs = make([]string, 0, len(params.Images))
		for i, img := range params.Images {
			item, err := imageToInputString(img)
			if err != nil {
				return nil, ValidationError{Op: ValidationOpImage, Err: fmt.Errorf("index %d: %w", i, err)}
			}
			inputs = append(inputs, item)
		}
	}

	model := strings.TrimSpace(params.Model)
	if model == "" {
		model = c.model
	}

	if strings.TrimSpace(model) == "" {
		return nil, ValidationError{Op: ValidationOpParams, Err: fmt.Errorf("missing model name")}
	}

	payload := embeddingRequest{
		Input: inputs,
		Model: model,
	}
	if params.Dimensions > 0 {
		payload.Dimensions = &params.Dimensions
	}
	if strings.TrimSpace(params.EncodingFormat) != "" {
		format := strings.ToLower(strings.TrimSpace(params.EncodingFormat))
		if format != "float" && format != "base64" {
			return nil, ValidationError{Op: ValidationOpParams, Err: fmt.Errorf("invalid encoding_format")}
		}
		payload.EncodingFormat = format
	}
	if strings.TrimSpace(params.User) != "" {
		payload.User = params.User
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
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, RequestError{Op: RequestOpCallService, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		if message := parseAPIError(respBody); message != "" {
			return nil, RequestError{
				Op:  "http status",
				Err: fmt.Errorf("%s: %s", resp.Status, message),
			}
		}
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

	results := make([]Embedding, 0, len(parsed.Data))
	for i, item := range parsed.Data {
		result, err := parseEmbedding(item.Embedding)
		if err != nil {
			return nil, ResponseError{
				Op:  ResponseOpResult,
				Err: fmt.Errorf("index %d: %w", i, err),
			}
		}
		results = append(results, result)
	}
	return results, nil
}

func imageToInputString(img Image) (string, error) {
	if strings.TrimSpace(img.URL) != "" {
		return img.URL, nil
	}
	if strings.TrimSpace(img.Base64) != "" {
		return img.Base64, nil
	}
	if img.Reader == nil {
		return "", ValidationError{Op: ValidationOpImage, Err: fmt.Errorf("missing image reader")}
	}

	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	if _, err := io.Copy(enc, img.Reader); err != nil {
		return "", RequestError{Op: RequestOpReadImageData, Err: err}
	}
	if err := enc.Close(); err != nil {
		return "", RequestError{Op: RequestOpEncodeImage, Err: err}
	}

	return buf.String(), nil
}

type apiErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func parseAPIError(payload []byte) string {
	if len(payload) == 0 {
		return ""
	}
	var parsed apiErrorResponse
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return ""
	}
	return strings.TrimSpace(parsed.Error.Message)
}

func parseEmbedding(raw json.RawMessage) (Embedding, error) {
	if len(raw) == 0 {
		return Embedding{}, fmt.Errorf("missing embedding")
	}
	var encoded string
	if err := json.Unmarshal(raw, &encoded); err == nil {
		if strings.TrimSpace(encoded) == "" {
			return Embedding{}, fmt.Errorf("empty embedding")
		}
		return Embedding{Base64: encoded}, nil
	}

	var floats64 []float64
	if err := json.Unmarshal(raw, &floats64); err == nil {
		if len(floats64) == 0 {
			return Embedding{}, fmt.Errorf("empty embedding")
		}
		vec := make([]float32, 0, len(floats64))
		for _, v := range floats64 {
			vec = append(vec, float32(v))
		}
		return Embedding{Vector: vec}, nil
	}

	return Embedding{}, fmt.Errorf("invalid embedding format")
}
