# Imbeddings Go SDK

Minimal Go client for the Imbeddings image-embeddings service (ViT-based models like DINOv3/DINOv2). It sends images as URLs or base64-encoded data and returns CLS and MEAN embeddings.

## Install

```bash
go get github.com/kotylevskiy/imbeddings-go-sdk
```

## Usage

```go
package main

import (
    "context"
    "log"
    "net/http"

    imbeddings "github.com/kotylevskiy/imbeddings-go-sdk"
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

    log.Printf("mean dims: %d, cls dims: %d", len(vectors[0].Mean), len(vectors[0].Cls))
}
```

## Notes

- Endpoint: `POST /v1/embeddings`
- Inputs: image URL or base64-encoded image data
- Outputs: CLS and MEAN embeddings (float arrays)
- The server must already be running and accessible.

## Configuration

If you want to align with environment-based configuration, the SDK exposes a helper struct:

- `ImbeddingsServiceOptions` (env vars: `IMBEDDINGS_BASE_URL`, `IMBEDDINGS_MODEL_ID`, `IMBEDDINGS_MAX_WIDTH`, `IMBEDDINGS_MAX_HEIGHT`)

