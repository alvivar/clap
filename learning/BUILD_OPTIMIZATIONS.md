# Go Build Optimizations

## Overview

This guide collects practical techniques for optimizing Go builds for size, performance, and reproducibility.
Always measure before and after changes.

## Baseline: default `go build`

By default, `go build`:

- Includes symbol tables and DWARF debugging info
- Enables standard compiler optimizations
- Balances size and speed
- Compiles all imported packages and dependencies

```bash
go build -o myapp main.go
```

## Quick start recipes

### Development (debug-friendly)

```bash
go build -o myapp main.go
```

### Balanced production (recommended)

```bash
go build -trimpath -ldflags="-s -w" -o myapp main.go
```

### Smallest binary

```bash
go build -ldflags="-s -w" -o myapp main.go
upx -9 myapp
```

### Fastest execution (PGO)

```bash
go build -pgo=auto -o myapp main.go
```

### Inspect compiler optimizations

```bash
go build -gcflags="all=-m" main.go
```

## Size optimization

### 1) Strip debug information

```bash
go build -ldflags="-s -w" -o myapp main.go
```

- `-s`: omit symbol table and debug info
- `-w`: omit DWARF debugging info

Trade-offs:

- ✅ significantly smaller binaries
- ✅ faster downloads and deployments
- ❌ cannot use debuggers (gdb, delve)
- ❌ stack traces may be less informative

### 2) Trim file paths (reproducibility + size)

```bash
go build -trimpath -ldflags="-s -w" -o myapp main.go
```

- Removes absolute paths from the binary
- Helps create reproducible builds

### 3) Build tags for conditional compilation

Use build tags to include optional code only when needed.

```go
// +build imageprocessing

package imageutils

import "image"

func ProcessImage(img *image.RGBA) {
    // Image processing functionality
}
```

```bash
# Include image processing
go build -tags=imageprocessing -o myapp main.go

# Exclude image processing (smaller binary)
go build -o myapp main.go
```

Use cases:

- Platform-specific features
- Optional functionality
- Development vs. production builds
- Feature flags

### 4) Dependency hygiene

```bash
go mod tidy
go mod graph
go build -o myapp main.go
go tool nm -size myapp | sort -rn | head -20
```

Tips:

- Avoid importing large packages for small features
- Review transitive dependencies regularly

### 5) UPX compression

```bash
go build -ldflags="-s -w" -o myapp main.go
upx -9 myapp
```

Trade-offs:

- ✅ extremely small binaries
- ✅ great for distribution
- ❌ slower startup (decompression overhead)
- ❌ some antivirus software may flag compressed binaries

### 6) Linker flags and metadata

```bash
go build -ldflags="-s -w -X main.version=1.0.0" -o myapp main.go
go build -ldflags="-linkmode external" -o myapp main.go
```

Useful `-ldflags`:

- `-X`: set a string variable at link time
- `-linkmode`: internal or external linking
- `-extldflags`: pass flags to an external linker

## Performance optimization

### 1) Profile-guided optimization (PGO)

PGO uses runtime profiles to guide compiler optimizations.

Collect a profile:

```go
import (
    "os"
    "runtime/pprof"
)

func main() {
    f, _ := os.Create("cpu.pprof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Your application code
}
```

Build with PGO:

```bash
./myapp
mv cpu.pprof default.pgo
go build -pgo=auto -o myapp-optimized main.go
```

Notes:

- Go 1.20 (preview), Go 1.21 (stable)
- Typical speedups: 5–20% on hot paths

### 2) Compiler optimization flags (`-gcflags`)

Common patterns:

```bash
go build -gcflags="all=-l" -o myapp main.go      # disable inlining
go build -gcflags="all=-B" -o myapp main.go      # disable bounds checks (unsafe)
go build -gcflags="all=-l=4" -o myapp main.go    # increase inlining budget
go build -gcflags="all=-m" -o myapp main.go      # print optimization decisions
```

Useful flags:

- `-m`: print optimization decisions
- `-l` / `-l=N`: control inlining
- `-N`: disable optimizations
- `-B`: disable bounds checks
- `-wb=false`: disable write barriers

### 3) Efficient code patterns

**String building**

```go
// ❌ Inefficient
var result string
for i := 0; i < 1000; i++ {
    result += "hello" // Creates new string each iteration
}

// ✅ Efficient
var sb strings.Builder
sb.Grow(5000) // Pre-allocate if size is known
for i := 0; i < 1000; i++ {
    sb.WriteString("hello")
}
result := sb.String()
```

**Memory pre-allocation**

```go
// ❌ Inefficient
var items []Item
for i := 0; i < 1000; i++ {
    items = append(items, Item{})
}

// ✅ Efficient
items := make([]Item, 0, 1000)
for i := 0; i < 1000; i++ {
    items = append(items, Item{})
}
```

**Avoiding allocations**

```go
// ❌ Creates allocations
func ProcessData(data []byte) []byte {
    return append([]byte("prefix: "), data...)
}

// ✅ Reuse buffers
var bufPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessData(data []byte) []byte {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()
    buf.WriteString("prefix: ")
    buf.Write(data)
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    bufPool.Put(buf)
    return result
}
```

## Diagnostics and measurement

### Benchmarks and profiling

```bash
go test -bench=. -benchmem
go test -bench=. -cpuprofile=cpu.out
go test -bench=. -memprofile=mem.out
go tool pprof cpu.out
go tool pprof mem.out
```

### Escape analysis

```bash
go build -gcflags="-m -m" main.go 2>&1 | grep "escapes to heap"
```

Example:

```go
// ❌ Escapes to heap
func CreatePoint() *Point {
    p := Point{X: 1, Y: 2}
    return &p // p escapes to heap
}

// ✅ Stack allocation
func CreatePoint(p *Point) {
    p.X = 1
    p.Y = 2
}

// Usage
var p Point
CreatePoint(&p)
```

### Binary size inspection

```bash
go build -o myapp main.go
go tool nm -size myapp | sort -rn | head -20
ls -lh myapp
du -h myapp
```

## Automation and reproducible builds

### Makefile example (tabs required)

```makefile
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -s -w -X main.version=$(VERSION)
GCFLAGS := all=-trimpath

.PHONY: build
build:
	go build -ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -o bin/myapp main.go

.PHONY: build-optimized
build-optimized:
	go build -ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -pgo=auto -o bin/myapp main.go
	upx -9 bin/myapp

.PHONY: profile
profile:
	go build -o bin/myapp-profile main.go
	./bin/myapp-profile
	mv cpu.pprof default.pgo
```

### Documenting builds

```go
// build.go
//go:build ignore

package main

// Build instructions:
// Production: go run build.go -prod
// Development: go run build.go -dev
// With PGO: go run build.go -prod -pgo
```

### Cross-platform builds

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o myapp-linux-amd64 main.go
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o myapp-windows-amd64.exe main.go
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o myapp-darwin-arm64 main.go
```

### Reproducibility tips

- Use `-trimpath` (or `-gcflags="all=-trimpath"`) to remove absolute paths.
- Pin your Go version and keep `go.mod` tidy.

## Development vs production builds

**Development**

```bash
go build -o myapp main.go
```

**Production**

```bash
go build -trimpath -ldflags="-s -w" -o myapp main.go
```

## Testing after optimization

1. Functional testing: ensure features still work.
2. Performance testing: benchmark key operations.
3. Load testing: test under realistic conditions.
4. Binary size: verify size reductions.
5. Startup time: check UPX impact.

## Reference tables

### Size reduction comparison

| Method             | Typical Reduction | Trade-offs            |
| ------------------ | ----------------- | --------------------- |
| Default build      | 0% (baseline)     | Full debug info       |
| `-ldflags="-s -w"` | 30-40%            | No debugging          |
| + Build tags       | 40-60%            | Reduced functionality |
| + UPX compression  | 50-70%            | Slower startup        |

### Performance improvement comparison

| Method                 | Typical Improvement | Complexity |
| ---------------------- | ------------------- | ---------- |
| PGO                    | 5-20%               | Medium     |
| Code optimization      | 10-100%+            | High       |
| Proper data structures | 10-1000%+           | High       |
| Compiler flags         | 0-5%                | Low        |

### Go version features

| Go Version | Notable Optimization Features          |
| ---------- | -------------------------------------- |
| 1.7+       | Better SSA compiler backend            |
| 1.11+      | Module support, reproducible builds    |
| 1.17+      | Register-based calling convention      |
| 1.18+      | Generics (can reduce code duplication) |
| 1.20+      | PGO preview, better inlining           |
| 1.21+      | PGO stable, ~15% faster                |
| 1.22+      | Better loop optimizations              |

## Resources

- https://pkg.go.dev/cmd/compile
- https://github.com/golang/go/wiki/Performance
- https://go.dev/doc/pgo
- https://pkg.go.dev/cmd/go#hdr-Build_modes
- https://go.dev/doc/effective_go
