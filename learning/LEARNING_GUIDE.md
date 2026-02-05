# Go Learning Guide (clap)

## Overview

This folder contains learning-oriented notes for Go, anchored around the `clap` CLI example in `main.go`.
Use it as a small learning path: project setup, language features, built-ins, and build optimization.

## Learning path

1. Project setup and structure → `docs/go-project-primer.md`
2. Language features used in this repo → “Applied Go features in `main.go`” (below)
3. Core built-ins → `docs/go-builtin-functions.md`
4. Shipping builds → `docs/go-build-optimizations.md`

## Applied Go features in `main.go`

### Program structure

- `package main` defines an executable program.
- Import blocks group standard library dependencies.
- `func main()` is the entry point.
- `func shouldPrintFile(...) bool` is a helper function returning a boolean.

### CLI flags and arguments

- `flag.String("o", "clap.txt", "output filename")` defines a flag and returns a `*string`.
- `flag.Parse()` parses flags.
- `flag.Args()` returns remaining positional arguments.
- `*outputFilename` demonstrates pointer dereference.

### Variables, slices, and iteration

- Short declarations with `:=` infer types.
- Slicing with `args[1:]`.
- `len(args)` checks slice length.
- `for _, filter := range filters` iterates over a slice.

### Control flow and error handling

- `if ... { ... }` guards required arguments and checks errors.
- Early returns use `return nil` and `return err`.
- `os.Exit(1)` exits with a non-zero status on failure.
- Multiple return values are captured as `value, err`.

### Filesystem and I/O

- `filepath.Join` and `filepath.Base` build and parse paths.
- `os.Create` and `os.ReadFile` handle file creation and reads.
- `bufio.NewWriter` buffers output for better performance.
- `defer` ensures cleanup with `Close` and `Flush`.

### Directory walking and callbacks

- `filepath.WalkDir` traverses the directory tree.
- An anonymous function (function literal) is used as the callback.
- `fs.DirEntry` and `d.Info()` expose file metadata.

### Strings and formatted output

- `strings.TrimSpace` and `strings.Contains` process filters.
- `fmt.Println`, `fmt.Printf`, and `fmt.Fprintf` produce formatted output.

## Practice ideas

- Add a `-ext` flag to filter by extension (e.g., `.go`, `.md`).
- Add a `-ignore` flag to skip directories.
- Replace `os.ReadFile` with streaming reads for large files.

## Related docs

- Project setup: `docs/go-project-primer.md`
- Built-ins: `docs/go-builtin-functions.md`
- Build optimization: `docs/go-build-optimizations.md`
