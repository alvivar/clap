package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outputFilename := flag.String("o", "clap.file", "output filename")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("👏 Clap slaps all your files into one!")
		fmt.Println("Usage: clap [-o filename] <path> [extensions...]")
		os.Exit(1)
	}

	path := args[0]
	extensions := normalizeExtensions(args[1:])

	var contentBuilder strings.Builder

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", filePath, err)
			return err
		}

		if info.IsDir() || !shouldPrintFile(filePath, extensions) {
			return nil
		}

		fmt.Printf("%s (%d bytes)\n", filePath, info.Size())

		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			return nil
		}

		contentBuilder.WriteString("=== ")
		contentBuilder.WriteString(filePath)
		contentBuilder.WriteString(" ===\n")
		contentBuilder.Write(content)
		contentBuilder.WriteString("\n\n")

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", path, err)
		os.Exit(1)
	}

	outputPath := filepath.Join(path, *outputFilename)
	if err := os.WriteFile(outputPath, []byte(contentBuilder.String()), 0644); err != nil {
		fmt.Printf("Error writing output file %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("Content written to %s\n", outputPath)
}

// normalizeExtensions converts extensions to a map with leading dots and lowercase.
// Returns nil if no extensions provided (accept all files).
func normalizeExtensions(extensions []string) map[string]bool {
	if len(extensions) == 0 {
		return nil
	}

	extMap := make(map[string]bool, len(extensions))
	for _, ext := range extensions {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		extMap[strings.ToLower(ext)] = true
	}
	return extMap
}

// shouldPrintFile returns true if the file matches the extension filter.
// If extensions is nil, all files are included.
func shouldPrintFile(filePath string, extensions map[string]bool) bool {
	if extensions == nil {
		return true
	}
	ext := strings.ToLower(filepath.Ext(filePath))
	return extensions[ext]
}
