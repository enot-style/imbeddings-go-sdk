# Imbeddings Go SDK

A lightweight, idiomatic Go client for accessing the [Imbeddings API](https://github.com/enot-style/imbeddings) — an image-embeddings service that produces dense vector embeddings from images (URLs, base64/data URIs). Typical use cases include building vector search, similarity ranking, RAG pipelines, or AI-powered image analysis workflows.

## Install

```bash
go get github.com/enot-style/imbeddings-go-sdk
```

## Usage

```go
package main

import (
    "context"
    "log"
    "net/http"

    imbeddings "github.com/enot-style/imbeddings-go-sdk"
)

func main() {
    client, err := imbeddings.NewClient("http://localhost:8000", "facebook/dinov3-vits16-pretrain-lvd1689m", http.DefaultClient)
    if err != nil {
        log.Fatal(err)
    }

    params := imbeddings.ImbeddingsParams{
        Images: []imbeddings.Image{
            imbeddings.ImageURL("https://example.com/cat.jpg"),
        },
    }

    vectors, err := client.New(context.Background(), params)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("dims: %d", len(vectors[0].Vector))
}
```

## Notes

- Endpoint: `POST /v1/embeddings`
- Inputs: image URL, data URI, or base64-encoded image data
- Outputs: a single embedding per image (float array or base64 string)
- The server must already be running and accessible.

Use `ImageBase64(...)` if you already have base64/data URI content.

You can also pass raw strings via `ImbeddingsParams.Inputs` if you don't need helpers.

## Options

`ImbeddingsParams` supports:
- `EncodingFormat`: `"float"` (default) or `"base64"`
- `Dimensions`: optional truncation size
- `User`: forwarded to the API for parity

## Configuration

Use `client.SetAPIKey(...)` if the service is configured with API keys.
