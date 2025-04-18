# go-asgardeo

Go SDK for WSO2 Asgardeo Management API.

## Installation

```bash
go get github.com/thilinashashimalsenarath/go-asgardeo/management
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/thilinashashimalsenarath/go-asgardeo/management"
)

func main() {
    // Initialize management client with client credentials
    client, err := management.New(
        "https://<your-tenant>.asgardeo.io",
        management.WithClientCredentials(context.Background(),
            "YOUR_CLIENT_ID",
            "YOUR_CLIENT_SECRET",
        ),
    )
    if err != nil {
        log.Fatalf("failed to initialize client: %v", err)
    }

    // List applications
    appsResp, err := client.Applications().List()
    if err != nil {
        log.Fatalf("failed to list applications: %v", err)
    }
    fmt.Printf("Total applications: %d\n", appsResp.TotalResults)
    for _, app := range appsResp.Applications {
        fmt.Printf("- %s (ID: %s)\n", app.Name, app.ID)
    }

    // Create a new application
    newApp := management.Application{
        Name:        "MyApp",
        Description: "Example application",
    }
    created, err := client.Applications().Create(newApp)
    if err != nil {
        log.Fatalf("failed to create application: %v", err)
    }
    fmt.Printf("Created application with ID: %s\n", created.ID)
```

## Examples

Run the application management sample:
```bash
go run examples/application-management-sample/main.go
```

## Integration Tests

The SDK includes an integration test against a real Asgardeo tenant. It requires the following environment variables:

```bash
export ASGARDEO_BASE_URL="https://<your-tenant>.asgardeo.io"
export ASGARDEO_CLIENT_ID="YOUR_CLIENT_ID"
export ASGARDEO_CLIENT_SECRET="YOUR_CLIENT_SECRET"
```

Run the test with the `integration` build tag. Flags must come _before_ the package path. By default Go only shows failures; to see passing tests also add `-v`:

```bash
# Run only integration tests (flags before path)
go test -tags=integration -v ./integration

# Or run all tests (unit + integration)
go test -timeout 60s -tags=integration -v ./...
```

## Services

- Applications
- (Future) User management
- (Future) Attribute management
- ...

## Contributing

Feel free to submit issues and pull requests to add more services and features.
