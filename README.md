# bitHuman Go SDK

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/bithuman-go/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/bithuman-go/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/bithuman-go/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/bithuman-go/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/bithuman-go/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/bithuman-go/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/bithuman-go
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/bithuman-go
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/bithuman-go
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fbithuman-go
 [loc-svg]: https://tokei.rs/b1/github/plexusone/bithuman-go
 [repo-url]: https://github.com/plexusone/bithuman-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/bithuman-go/blob/main/LICENSE

Go SDK for the [bitHuman](https://www.bithuman.ai/) Real-Time Avatar Animation API.

bitHuman creates digital avatars that lip-sync to audio in real time. Audio in, animated video out at 25 FPS.

## Installation

```bash
go get github.com/plexusone/bithuman-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/plexusone/bithuman-go"
    "github.com/plexusone/bithuman-go/api"
)

func main() {
    // Create client (uses BITHUMAN_API_KEY env var if not specified)
    client, err := bithuman.NewClient(
        bithuman.WithAPIKey(os.Getenv("BITHUMAN_API_KEY")),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Validate credentials
    resp, err := client.Validate(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Credentials valid: %v\n", resp.Valid)

    // List agents
    agents, err := client.Agents().List(ctx, api.ListAgentsParams{})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d agents\n", len(agents.Agents))

    // Create a real-time session
    session, err := client.Sessions().Create(ctx, &api.CreateSessionRequest{
        AgentID: "your-agent-id",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Session created: %s\n", session.ID)
    fmt.Printf("LiveKit URL: %s\n", session.LivekitURL.Value)
}
```

## Features

- 🤖 **Agents**: Create, manage, and interact with avatar agents
- 📡 **Sessions**: Real-time conversation sessions with LiveKit integration
- 🗣️ **TTS**: Text-to-speech synthesis with 30+ language support
- 🎬 **Videos**: Generate talking video MP4 files
- 📁 **Files**: Upload images, video, audio, and documents
- 🔔 **Webhooks**: Event notifications for async operations

## Environment Variables

| Variable | Description |
|----------|-------------|
| `BITHUMAN_API_KEY` | Your bitHuman API key (from Developer → API Keys) |

## API Reference

### Client Options

```go
client, _ := bithuman.NewClient(
    bithuman.WithAPIKey("your-api-key"),
    bithuman.WithBaseURL("https://api.bithuman.ai"),  // Optional
    bithuman.WithTimeout(60 * time.Second),           // Optional
    bithuman.WithHTTPClient(customClient),            // Optional
)
```

### Services

| Service | Description |
|---------|-------------|
| `client.Agents()` | Avatar agent management |
| `client.Sessions()` | Real-time conversation sessions |
| `client.TTS()` | Text-to-speech synthesis |
| `client.Videos()` | Talking video generation |
| `client.Files()` | File uploads |
| `client.Billing()` | Account balance |
| `client.Webhooks()` | Event notifications |

### Low-Level API Access

For endpoints not covered by the high-level services:

```go
apiClient := client.API()
// Use ogen-generated methods directly
```

## LiveKit Integration

bitHuman sessions can connect to LiveKit for real-time WebRTC streaming:

```go
session, _ := client.Sessions().Create(ctx, &api.CreateSessionRequest{
    AgentID:      "your-agent-id",
    LivekitURL:   api.NewOptString("wss://your-livekit-server"),  // Optional external LiveKit
    LivekitToken: api.NewOptString("your-token"),                  // Optional external token
})

// Use session.LivekitURL and session.LivekitToken to connect
```

## Development

### Regenerating the API Client

The SDK uses [ogen](https://github.com/ogen-go/ogen) to generate Go code from the OpenAPI specification:

```bash
# Install ogen
go install github.com/ogen-go/ogen/cmd/ogen@latest

# Regenerate API client
./generate.sh
```

### Running Tests

```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```

## Related

- [bitHuman Documentation](https://docs.bithuman.ai/)
- [bitHuman API Reference](https://docs.bithuman.ai/api-reference/overview)
- [LiveKit](https://livekit.io/) - Real-time WebRTC platform

## License

MIT License - see [LICENSE](LICENSE) for details.
