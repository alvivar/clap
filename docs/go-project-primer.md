# Go Project Primer Guide

A complete guide to starting and structuring a project with Go (Golang).

## Table of Contents

-   [Prerequisites](#prerequisites)
-   [Initial Setup](#initial-setup)
-   [Project Structure](#project-structure)
-   [Managing Dependencies](#managing-dependencies)
-   [Building and Running](#building-and-running)
-   [Best Practices](#best-practices)

---

## Prerequisites

### 1. Install Go

Download and install Go from the official website: [https://go.dev/dl/](https://go.dev/dl/)

Verify your installation:

```bash
go version
```

You should see output like: `go version go1.21.x windows/amd64`

---

## Initial Setup

### 2. Create Your Project Directory

```bash
mkdir my-go-project
cd my-go-project
```

### 3. Initialize a Go Module

Go modules are the standard way to manage dependencies and versions:

```bash
go mod init github.com/yourusername/my-go-project
```

**Tips:**

-   Replace `github.com/yourusername/my-go-project` with your actual module path
-   For personal projects, you can use any name: `go mod init my-project`
-   This creates a `go.mod` file that tracks dependencies

### 4. Create Your First Go File

Create `main.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

**Key Points:**

-   `package main` declares this as an executable program
-   `main()` function is the entry point
-   Import statements bring in standard library or external packages

---

## Project Structure

### Basic Structure

For small projects:

```
my-go-project/
├── go.mod
├── go.sum          (created after adding dependencies)
├── main.go
└── README.md
```

### Organized Structure

As your project grows, organize into packages:

```
my-go-project/
├── go.mod
├── go.sum
├── main.go
├── cmd/            (command-line applications)
│   └── server/
│       └── main.go
├── internal/       (private application code)
│   ├── handlers/
│   │   └── handler.go
│   └── models/
│       └── user.go
├── pkg/            (public library code)
│   └── utils/
│       └── helpers.go
└── README.md
```

**Directory Conventions:**

-   `cmd/` - Main applications for the project
-   `internal/` - Private code (cannot be imported by external projects)
-   `pkg/` - Public library code (can be imported)
-   `api/` - API definitions (OpenAPI, Protocol Buffers)
-   `web/` - Web assets (templates, static files)
-   `test/` - Additional test files and test data

### Example: Multi-Package Project

**greetings/greetings.go**:

```go
package greetings

import "fmt"

func Hello(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
```

**main.go**:

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

---

## Managing Dependencies

### Adding Dependencies

Use `go get` to add external packages:

```bash
go get rsc.io/quote
go get github.com/gorilla/mux@v1.8.0  # specific version
```

### Using Dependencies

Import and use in your code:

```go
package main

import (
    "fmt"
    "rsc.io/quote"
)

func main() {
    fmt.Println(quote.Go())
}
```

### Tidy Dependencies

Remove unused dependencies and add missing ones:

```bash
go mod tidy
```

### View Dependencies

```bash
go list -m all        # List all dependencies
go mod graph          # Show dependency graph
```

---

## Building and Running

### Run Without Building

```bash
go run main.go
```

For multi-file packages:

```bash
go run .
```

### Build an Executable

```bash
go build
```

This creates an executable (e.g., `my-go-project.exe` on Windows).

Run it:

```bash
./my-go-project       # Linux/Mac
my-go-project.exe     # Windows
```

### Build with Custom Name

```bash
go build -o myapp
```

### Cross-Platform Builds

Build for different platforms:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o myapp-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o myapp-mac

# Windows
GOOS=windows GOARCH=amd64 go build -o myapp.exe
```

### Testing

Run tests:

```bash
go test ./...         # Test all packages
go test -v ./...      # Verbose output
go test -cover ./...  # With coverage
```

---

## Best Practices

### 1. Code Organization

-   Keep packages focused on a single responsibility
-   Use descriptive package names (avoid `util`, `common`, `helpers`)
-   Organize code by functionality, not by type

### 2. Naming Conventions

-   Use `camelCase` for private functions/variables
-   Use `PascalCase` for public (exported) functions/variables
-   Keep names short but meaningful
-   Avoid stuttering: `user.UserID` → `user.ID`

### 3. Error Handling

```go
result, err := someFunction()
if err != nil {
    // Handle error
    return err
}
// Use result
```

### 4. Go Formatting

Always format your code:

```bash
go fmt ./...
```

Better yet, use `gofmt` or tools like `goimports`:

```bash
goimports -w .
```

### 5. Version Control

Initialize Git:

```bash
git init
git add .
git commit -m "Initial commit"
```

Create `.gitignore`:

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

### 6. Documentation

-   Write package documentation above the package declaration
-   Document exported functions, types, and constants
-   Use `godoc` format

```go
// Package greetings provides functions for greeting users.
package greetings

// Hello returns a greeting message for the given name.
func Hello(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
```

### 7. Common Commands Cheat Sheet

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

---

## Quick Start Template

Here's a quick template to start any Go project:

```bash
# 1. Create and navigate to project directory
mkdir my-project && cd my-project

# 2. Initialize module
go mod init github.com/yourusername/my-project

# 3. Create main.go
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
EOF

# 4. Run it
go run .

# 5. Initialize git
git init
echo "*.exe" > .gitignore
git add .
git commit -m "Initial commit"
```

---

## Resources

-   **Official Documentation**: [https://go.dev/doc/](https://go.dev/doc/)
-   **Go Tour**: [https://go.dev/tour/](https://go.dev/tour/)
-   **Effective Go**: [https://go.dev/doc/effective_go](https://go.dev/doc/effective_go)
-   **Go by Example**: [https://gobyexample.com/](https://gobyexample.com/)
-   **Standard Library**: [https://pkg.go.dev/std](https://pkg.go.dev/std)

---

## What's Next?

After completing this primer:

1. ✅ Complete the Go Tour for language fundamentals
2. ✅ Read Effective Go for idiomatic patterns
3. ✅ Build a small CLI tool or web API
4. ✅ Learn about testing with the `testing` package
5. ✅ Explore popular frameworks (Gin, Echo, Chi)
6. ✅ Study concurrency with goroutines and channels

Happy coding! 🚀
