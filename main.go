package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// fatalf is a tiny helper for fatal CLI errors.
//
// - format string works like fmt.Printf
// - args ...any is a variadic parameter (zero or more values)
// - os.Exit(1) ends the process with a non-zero status code
func fatalf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}

func main() {
	// flag.String returns a *string (pointer to string).
	// We dereference it later with *outputFilename after flag.Parse().
	outputFilename := flag.String("o", "clap.txt", "output filename")

	// You can override default help text by assigning flag.Usage.
	flag.Usage = func() {
		fmt.Println("👏 Clap slaps all your files into one!")
		fmt.Println("Usage: clap [-o filename] <path> [ext <extensions...>] [sub <substrings...>]")
		fmt.Println("  -o filename  output filename (default: clap.txt in current working directory)")
		fmt.Println("Filters are case-insensitive. Extensions may include or omit the leading dot.")
	}

	// Parse CLI flags and keep the remaining positional args in flag.Args().
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		fatalf("Error: path is required")
	}

	rootPath := args[0]
	extFilters, subFilters, err := parseFilterArgs(args[1:])
	if err != nil {
		fatalf("Error parsing filters: %v", err)
	}

	// Normalize to absolute paths for robust comparisons/IO.
	outputPath, err := filepath.Abs(*outputFilename)
	if err != nil {
		fatalf("Error resolving output path: %v", err)
	}

	rootPath, err = filepath.Abs(rootPath)
	if err != nil {
		fatalf("Error resolving path: %v", err)
	}

	// os.Create truncates existing files or creates a new one.
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fatalf("Error creating output file: %v", err)
	}
	// defer schedules this call for function exit (even on later returns).
	defer outputFile.Close()

	// bufio.NewWriter batches writes for better performance.
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	wroteAny := false

	// filepath.WalkDir recursively visits all entries under rootPath.
	// The callback is a function literal (anonymous function) that closes over
	// outer variables (rootPath/outputPath/writer/filters/wroteAny).
	err = filepath.WalkDir(rootPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", filePath, err)
			return err
		}

		// Skip directories, the output file itself, and non-matching files.
		if d.IsDir() || filePath == outputPath || !shouldPrintFile(filePath, extFilters, subFilters) {
			return nil
		}

		// Read whole file into memory (simple and fine for normal text files).
		content, err := os.ReadFile(filePath)
		if err != nil {
			// Non-fatal for a single file: log and continue walking.
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			return nil
		}

		// Convert to a path relative to scanned root for cleaner output.
		displayPath := filePath
		if relPath, relErr := filepath.Rel(rootPath, filePath); relErr == nil {
			if relPath == "." {
				displayPath = filepath.Base(filePath)
			} else {
				displayPath = relPath
			}
		}

		fmt.Printf("%s (%d bytes)\n", displayPath, len(content))

		// Add a blank separator only between sections (not at file end).
		if wroteAny {
			if _, err := writer.WriteString("\n\n"); err != nil {
				return err
			}
		}
		wroteAny = true

		if _, err := writer.WriteString("=== " + displayPath + " ===\n"); err != nil {
			return err
		}
		if _, err := writer.Write(content); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fatalf("Error walking the path %s: %v", rootPath, err)
	}

	fmt.Printf("Content written to %s\n", outputPath)
}

// parseFilterArgs parses CLI fragments like:
//
//	ext go md sub test config
//
// It returns:
//   - extension filters (without leading dot, lowercase)
//   - substring filters (lowercase)
//   - error for malformed input
func parseFilterArgs(args []string) ([]string, []string, error) {
	// Slices are dynamic views over arrays.
	var extFilters, subFilters []string

	// state holds current section name: "ext", "sub", or "" (none yet).
	state := ""
	hasValue := false

	for _, arg := range args { // range loops over slice elements
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}

		if arg == "ext" || arg == "sub" {
			// If previous section had no values, it's invalid.
			if state != "" && !hasValue {
				return nil, nil, fmt.Errorf("%s requires at least one value", state)
			}
			state, hasValue = arg, false
			continue
		}

		if state == "" {
			return nil, nil, fmt.Errorf("unexpected filter value %q; expected \"ext\" or \"sub\"", arg)
		}

		if state == "ext" {
			// Short variable declaration in if statement scope.
			if ext := strings.TrimPrefix(strings.ToLower(arg), "."); ext != "" {
				extFilters = append(extFilters, ext)
				hasValue = true
			}
			continue
		}

		subFilters = append(subFilters, strings.ToLower(arg))
		hasValue = true
	}

	if state != "" && !hasValue {
		return nil, nil, fmt.Errorf("%s requires at least one value", state)
	}

	return extFilters, subFilters, nil
}

// shouldPrintFile checks whether a file matches both filters (if provided).
//
// extFilters: exact extension match (case-insensitive, without dot)
// subFilters: filename contains at least one substring (case-insensitive)
func shouldPrintFile(filePath string, extFilters, subFilters []string) bool {
	fileName := filepath.Base(filePath)

	if len(extFilters) > 0 {
		fileExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")
		if !slices.Contains(extFilters, fileExt) {
			return false
		}
	}

	if len(subFilters) > 0 {
		lowerName := strings.ToLower(fileName)
		if !slices.ContainsFunc(subFilters, func(filter string) bool {
			return strings.Contains(lowerName, filter)
		}) {
			return false
		}
	}

	return true
}
