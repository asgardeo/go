# Go Asgardeo SDK

Typed Go client SDK for the Asgardeo Management API, enabling easy integration with Asgardeo services.

## Requirements

- Go 1.22 or later

## Installation

Use `go get` to install the SDK:

```bash
go get github.com/asgardeo/go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/asgardeo/go/pkg/config"
    "github.com/asgardeo/go/pkg/sdk"
)

func main() {
    // Configure the client with client credentials grant
    cfg := config.DefaultClientConfig().
        WithBaseURL("https://api.asgardeo.io/t/<tenant-domain>").
        WithTimeout(10 * time.Second).
        WithClientCredentials(
            "YOUR_CLIENT_ID",
            "YOUR_CLIENT_SECRET",
        )

    // Initialize the SDK client
    client, err := sdk.NewClient(cfg)
    if err != nil {
        log.Fatalf("failed to initialize SDK client: %v", err)
    }

    // List applications
    ctx := context.Background()
    resp, err := client.ApplicationClient.List(ctx, 10, 0)
    if err != nil {
        log.Fatalf("failed to list applications: %v", err)
    }

    fmt.Printf("Total applications: %d\n", *resp.TotalResults)
    for _, app := range *resp.Applications {
        fmt.Printf("- %s (ID: %s)\n", *app.Name, *app.Id)
    }
}
```

## Examples

A runnable example is available in the `examples/application` directory:

```bash
go run examples/application/main.go
```

## Services

- Applications

## Contributing

Contributions, issues, and feature requests are welcome! Feel free to open a pull request or an issue to get started.
