# Go Project Primer

## Overview

A practical guide to starting, structuring, and maintaining a Go project.

## Prerequisites

1. Install Go: https://go.dev/dl/
2. Verify the installation:
    ```bash
    go version
    ```
    Example output: `go version go1.21.x windows/amd64`

## Quick start (minimal project)

1. Create and enter a project directory:
    ```bash
    mkdir my-go-project
    cd my-go-project
    ```
2. Initialize a module:
    ```bash
    go mod init github.com/yourusername/my-go-project
    ```
    Tips:
    - Use your real module path for public projects.
    - For personal projects, a short name is fine (e.g., `go mod init my-project`).
3. Create `main.go`:

    ```go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello, World!")
    }
    ```

4. Run the program:
    ```bash
    go run .
    ```
5. Build an executable:
    ```bash
    go build -o myapp
    ```

## Project structure

### Small projects

```
my-go-project/
├── go.mod
├── go.sum
├── main.go
└── README.md
```

### Larger projects (conventional layout)

```
my-go-project/
├── cmd/                (command-line applications)
│   └── server/
│       └── main.go
├── internal/           (private application code)
│   ├── handlers/
│   │   └── handler.go
│   └── models/
│       └── user.go
├── pkg/                (public library code)
│   └── utils/
│       └── helpers.go
├── api/                (OpenAPI, Protobuf, etc.)
├── web/                (web assets/templates)
├── test/               (test data and helpers)
├── go.mod
├── go.sum
└── README.md
```

### Directory conventions

- `cmd/` – entry points for binaries.
- `internal/` – private packages (not importable outside the module).
- `pkg/` – reusable libraries (exported).
- `api/`, `web/`, `test/` – optional, based on project needs.

### Example: multi-package

**greetings/greetings.go**

```go
package greetings

import "fmt"

func Hello(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
```

**main.go**

```go
package main

import (
    "fmt"
    "my-go-project/greetings"
)

func main() {
    message := greetings.Hello("World")
    fmt.Println(message)
}
```

## Modules and dependencies

### Add dependencies

```bash
go get rsc.io/quote
go get github.com/gorilla/mux@v1.8.0
```

### Clean and inspect dependencies

```bash
go mod tidy
go list -m all
go mod graph
```

## Build, run, and test

### Run

```bash
go run main.go
go run .
```

### Build

```bash
go build
go build -o myapp
```

### Cross-platform builds

```bash
GOOS=linux GOARCH=amd64 go build -o myapp-linux
GOOS=darwin GOARCH=amd64 go build -o myapp-mac
GOOS=windows GOARCH=amd64 go build -o myapp.exe
```

### Test

```bash
go test ./...
go test -v ./...
go test -cover ./...
```

## Tooling and formatting

```bash
go fmt ./...
goimports -w .
go vet ./...
```

## Best practices

### Code organization

- Keep packages focused on a single responsibility.
- Organize code by feature or domain, not by type.

### Naming

- Use `camelCase` for unexported identifiers.
- Use `PascalCase` for exported identifiers.
- Avoid stuttering: `user.UserID` → `user.ID`.

### Error handling

```go
result, err := someFunction()
if err != nil {
    return err
}
```

### Documentation

- Write package docs above the `package` declaration.
- Document exported symbols using `godoc` style.

```go
// Package greetings provides functions for greeting users.
package greetings

// Hello returns a greeting message for the given name.
func Hello(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
```

### Version control

```bash
git init
git add .
git commit -m "Initial commit"
```

**.gitignore example**

```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test files
*.test
*.out

# Go workspace file
go.work

# IDE
.vscode/
.idea/
```

## Command cheat sheet

```bash
go run .              # Run the project
go build              # Compile the project
go test ./...         # Run all tests
go mod tidy           # Clean up dependencies
go fmt ./...          # Format all files
go vet ./...          # Check for common mistakes
go get -u             # Update all dependencies
go clean              # Remove build artifacts
```

## Quick start script (copy/paste)

```bash
mkdir my-project && cd my-project
go mod init github.com/yourusername/my-project

cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
EOF

go run .
git init
echo "*.exe" > .gitignore
git add .
git commit -m "Initial commit"
```

## Resources

- Official docs: https://go.dev/doc/
- Go Tour: https://go.dev/tour/
- Effective Go: https://go.dev/doc/effective_go
- Go by Example: https://gobyexample.com/
- Standard Library: https://pkg.go.dev/std

## Next steps

1. Complete the Go Tour for fundamentals.
2. Read Effective Go for idiomatic patterns.
3. Build a small CLI or web API.
4. Learn testing with the `testing` package.
5. Explore popular frameworks (Gin, Echo, Chi).
6. Study concurrency with goroutines and channels.
