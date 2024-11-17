# URL Shortener

A simple URL shortener service built with Go using only the standard library.

## Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   └── handlers.go
│   └── middleware/
│       └── middleware.go
├── Makefile
├── go.mod
└── README.md
```

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/AjayPoshak/url-shortener.git
cd url-shortener
```

2. Install dependencies:

```bash
go mod tidy
```

## Running the Server

There are several ways to run the server:

### Direct Go Run

```bash
go run cmd/server/main.go
```

### Using Make

```bash
# Run the server
make run

# Build the binary
make build

# Clean build artifacts
make clean
```

### Using Docker

```bash
podman compose watch
```

The server will start on `http://localhost:8095`

## Testing the API

You can test the endpoints using curl:

```bash
# Test home endpoint
curl http://localhost:8095/

# Test health endpoint
curl http://localhost:8095/health
```

## Development

### Project Layout

- `cmd/server/main.go` - Main application entry point
- `internal/handlers/` - HTTP request handlers
- `internal/middleware/` - HTTP middleware (logging, etc.)

### Building for Production

To build a production binary:

```bash
go build -ldflags="-w -s" -o url-shortener cmd/server/main.go
```

The binary will be created in the current directory and can be run with:

```bash
./url-shortener
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/awesome-feature`)
3. Commit your changes (`git commit -m 'Add awesome feature'`)
4. Push to the branch (`git push origin feature/awesome-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
