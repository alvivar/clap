# Go Build Optimizations Primer

A comprehensive guide to creating optimized Go builds for performance and size.

## Table of Contents

-   [Default Build Configuration](#default-build-configuration)
-   [Binary Size Optimization](#binary-size-optimization)
-   [Performance Optimization](#performance-optimization)
-   [Advanced Techniques](#advanced-techniques)
-   [Best Practices](#best-practices)
-   [Quick Reference](#quick-reference)

## Default Build Configuration

By default, when you run `go build`, the Go compiler:

-   **Includes debugging information**: Symbol tables and DWARF debugging data are included to facilitate development and debugging
-   **Enables optimizations**: The compiler applies standard optimizations without requiring flags
-   **Balances performance and size**: Aims for a reasonable middle ground between binary size and execution speed
-   **Includes all code**: Compiles all imported packages and their dependencies

Default build command:

```bash
go build -o myapp main.go
```

This produces a fully functional binary with debugging support, suitable for development but not optimized for production deployment.

## Binary Size Optimization

### 1. Strip Debug Information

The most effective way to reduce binary size is to strip symbol tables and debugging information.

```bash
go build -ldflags="-s -w" -o myapp main.go
```

**Flags explained:**

-   `-s`: Omits the symbol table and debug information
-   `-w`: Omits DWARF debugging information

**Impact:** Can reduce binary size by 30-40% depending on the codebase.

**Trade-offs:**

-   ✅ Significantly smaller binaries
-   ✅ Faster downloads and deployments
-   ❌ Cannot use debuggers (gdb, delve) with the binary
-   ❌ Stack traces may be less informative

### 2. Build Tags for Conditional Compilation

Use build tags to exclude unnecessary code and dependencies from compilation.

**Example:**

```go
// +build imageprocessing

package imageutils

import "image"

func ProcessImage(img *image.RGBA) {
    // Image processing functionality
}
```

Build with specific tags:

```bash
# Include image processing
go build -tags=imageprocessing -o myapp main.go

# Exclude image processing (smaller binary)
go build -o myapp main.go
```

**Use cases:**

-   Platform-specific features
-   Optional functionality
-   Development vs. production builds
-   Feature flags

### 3. UPX Compression

UPX (Ultimate Packer for Executables) can compress binaries even further.

```bash
# Build with stripped symbols
go build -ldflags="-s -w" -o myapp main.go

# Compress with UPX
upx -9 myapp
```

**Impact:** Can achieve 50-70% total size reduction.

**Trade-offs:**

-   ✅ Extremely small binaries
-   ✅ Great for distribution
-   ❌ Increased startup time (decompression overhead)
-   ❌ Some antivirus software may flag compressed binaries
-   ❌ Cannot easily analyze the binary

### 4. Dependency Management

Keep your dependencies lean:

```bash
# View dependencies
go mod graph

# Remove unused dependencies
go mod tidy

# Analyze what's contributing to binary size
go build -o myapp main.go
go tool nm -size myapp | sort -rn | head -20
```

**Tips:**

-   Avoid importing large packages for small features
-   Use interfaces to reduce coupling
-   Consider vendoring only what you need
-   Review transitive dependencies regularly

## Performance Optimization

### 1. Profile-Guided Optimization (PGO)

PGO uses runtime profiling data to guide compiler optimizations, tailoring builds to actual usage patterns.

**Step 1: Collect a profile**

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

**Step 2: Build with PGO**

```bash
# Run your app to generate profile
./myapp

# Rename profile for PGO
mv cpu.pprof default.pgo

# Build with PGO
go build -pgo=auto -o myapp-optimized main.go
```

**Impact:**

-   5-20% performance improvement on hot paths
-   Better inlining decisions
-   Optimized for real-world usage patterns

**Available since:** Go 1.20 (preview), Go 1.21 (stable)

### 2. Compiler Optimization Flags

Control compiler optimizations with `-gcflags`:

```bash
# Disable inlining (reduces size, may hurt performance)
go build -gcflags="all=-l" -o myapp main.go

# Disable bounds checking (faster, but unsafe)
go build -gcflags="all=-B" -o myapp main.go

# Increase inlining budget (better performance, larger binary)
go build -gcflags="all=-l=4" -o myapp main.go

# View all optimization decisions
go build -gcflags="all=-m" -o myapp main.go
```

**Common flags:**

-   `-m`: Print optimization decisions
-   `-l`: Control inlining (`-l` disables, `-l=N` sets budget)
-   `-N`: Disable optimizations
-   `-B`: Disable bounds checking
-   `-wb=false`: Disable write barriers

### 3. Linker Optimization Flags

Beyond `-s` and `-w`, other linker flags can help:

```bash
# Reduce binary size and improve performance
go build -ldflags="-s -w -X main.version=1.0.0" -o myapp main.go

# External linking (can enable additional optimizations)
go build -ldflags="-linkmode external" -o myapp main.go
```

**Useful ldflags:**

-   `-X`: Set string variable at link time
-   `-linkmode`: Control linking mode (internal/external)
-   `-extldflags`: Pass flags to external linker

### 4. Efficient Code Patterns

**String Building:**

```go
// ❌ Inefficient
var result string
for i := 0; i < 1000; i++ {
    result += "hello"  // Creates new string each iteration
}

// ✅ Efficient
var sb strings.Builder
sb.Grow(5000)  // Pre-allocate if size is known
for i := 0; i < 1000; i++ {
    sb.WriteString("hello")
}
result := sb.String()
```

**Memory Pre-allocation:**

```go
// ❌ Inefficient
var items []Item
for i := 0; i < 1000; i++ {
    items = append(items, Item{})  // May trigger multiple reallocations
}

// ✅ Efficient
items := make([]Item, 0, 1000)  // Pre-allocate capacity
for i := 0; i < 1000; i++ {
    items = append(items, Item{})
}
```

**Avoiding Allocations:**

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

## Advanced Techniques

### 1. Custom Build Scripts

Create a Makefile or build script for consistent builds:

```makefile
# Makefile
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

### 2. Cross-Platform Optimization

Build for multiple platforms with optimizations:

```bash
# Linux AMD64 with all optimizations
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o myapp-linux-amd64 main.go

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o myapp-windows-amd64.exe main.go

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o myapp-darwin-arm64 main.go
```

### 3. Escape Analysis

Understand and optimize memory allocations:

```bash
# See where allocations escape to heap
go build -gcflags="-m -m" main.go 2>&1 | grep "escapes to heap"
```

**Example optimization:**

```go
// ❌ Escapes to heap
func CreatePoint() *Point {
    p := Point{X: 1, Y: 2}
    return &p  // p escapes to heap
}

// ✅ Stack allocation
func CreatePoint(p *Point) {
    p.X = 1
    p.Y = 2
}

// Usage
var p Point
CreatePoint(&p)  // p stays on stack
```

### 4. Benchmarking and Profiling

Always measure before and after optimization:

```bash
# Run benchmarks
go test -bench=. -benchmem

# Generate CPU profile
go test -bench=. -cpuprofile=cpu.out

# Generate memory profile
go test -bench=. -memprofile=mem.out

# Analyze with pprof
go tool pprof cpu.out
go tool pprof mem.out
```

## Best Practices

### Development vs. Production Builds

**Development:**

```bash
# Keep debug info, fast compile times
go build -o myapp main.go
```

**Production:**

```bash
# Optimize for size and performance
go build -ldflags="-s -w" -gcflags="all=-trimpath" -o myapp main.go
```

### Build Reproducibility

Ensure consistent builds:

```bash
# Use -trimpath to remove absolute paths
go build -trimpath -ldflags="-s -w" -o myapp main.go

# Pin Go version
go mod tidy
```

### Testing After Optimization

Always verify your optimizations:

1. **Functional testing**: Ensure all features still work
2. **Performance testing**: Benchmark key operations
3. **Load testing**: Test under realistic conditions
4. **Binary size**: Verify size reductions
5. **Startup time**: Check if UPX compression affects startup

### Documentation

Document your build process:

```go
// build.go
//go:build ignore

package main

// Build instructions:
// Production: go run build.go -prod
// Development: go run build.go -dev
// With PGO: go run build.go -prod -pgo
```

## Quick Reference

### Common Build Commands

```bash
# Smallest binary
go build -ldflags="-s -w" -o myapp main.go && upx -9 myapp

# Fastest execution
go build -pgo=auto -o myapp main.go

# Balanced (recommended for production)
go build -trimpath -ldflags="-s -w" -o myapp main.go

# Debug build
go build -gcflags="all=-N -l" -o myapp main.go

# View compiler optimizations
go build -gcflags="all=-m" main.go

# Check binary size
ls -lh myapp
du -h myapp
```

### Size Reduction Comparison

| Method             | Typical Reduction | Trade-offs            |
| ------------------ | ----------------- | --------------------- |
| Default build      | 0% (baseline)     | Full debug info       |
| `-ldflags="-s -w"` | 30-40%            | No debugging          |
| + Build tags       | 40-60%            | Reduced functionality |
| + UPX compression  | 50-70%            | Slower startup        |

### Performance Improvement Comparison

| Method                 | Typical Improvement | Complexity |
| ---------------------- | ------------------- | ---------- |
| PGO                    | 5-20%               | Medium     |
| Code optimization      | 10-100%+            | High       |
| Proper data structures | 10-1000%+           | High       |
| Compiler flags         | 0-5%                | Low        |

### Go Version Features

| Go Version | Notable Optimization Features          |
| ---------- | -------------------------------------- |
| 1.7+       | Better SSA compiler backend            |
| 1.11+      | Module support, reproducible builds    |
| 1.17+      | Register-based calling convention      |
| 1.18+      | Generics (can reduce code duplication) |
| 1.20+      | PGO preview, better inlining           |
| 1.21+      | PGO stable, ~15% faster                |
| 1.22+      | Better loop optimizations              |

## Conclusion

Optimizing Go builds requires balancing multiple factors:

-   **For production deployment**: Use `-ldflags="-s -w"` and `-trimpath`
-   **For maximum performance**: Use PGO with production profiles
-   **For minimum size**: Combine stripping, build tags, and UPX
-   **For development**: Use default builds with debug information

Always measure the impact of optimizations on your specific application and use cases. Profile before optimizing, and verify that optimizations provide real benefits without introducing bugs or unacceptable trade-offs.

## Additional Resources

-   [Official Go Compiler Documentation](https://pkg.go.dev/cmd/compile)
-   [Go Performance Wiki](https://github.com/golang/go/wiki/Performance)
-   [Profile-Guided Optimization](https://go.dev/doc/pgo)
-   [Go Build Modes](https://pkg.go.dev/cmd/go#hdr-Build_modes)
-   [Effective Go](https://go.dev/doc/effective_go)
