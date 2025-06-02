# Memcached Go

Memcached Go is a lightweight, in-memory key-value store written in Go, inspired by the original Memcached project. It implements the core Memcached protocol, allowing clients to store, retrieve, and delete data efficiently.

### Goal

This project aims to provide a simple, idiomatic Go implementation of the Memcached protocol, suitable for experimentation, learning, and lightweight use cases.

## Project Structure

```
memcached-go
├── cmd
│   └── main.go        # Entry point of the application
├── internal
│   └── server
│       └── server.go       # Server commands
│   └── lexer
│       └── lexer.go       # Lexigraphically Parse Memcached Commands
├── go.mod             # Module definition
└── go.sum             # Dependency checksums
```

## Getting Started

To get started with the application, follow these steps:

1. Clone the repository:
   ```
   git clone <repository-url>
   ```

2. Navigate to the project directory:
   ```
   cd memcached-go
   ```

3. Install the dependencies:
   ```
   go mod tidy
   ```

4. Run the application:
   ```
   go run cmd/main.go
   ```

## License

This project is licensed under the MIT License. See the LICENSE file for more details.