package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outputFilename := flag.String("o", "clap.txt", "output filename")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("👏 Clap slaps all your files into one!")
		fmt.Println("Usage: clap [-o filename] <path> [extensions...]")
		fmt.Println("Error: path is required")
		os.Exit(1)
	}

	path := args[0]
	extensions := normalizeExtensions(args[1:])
	outputPath := filepath.Join(path, *outputFilename)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Use buffered writer for better performance
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	err = filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", filePath, err)
			return err
		}

		// Skip directories, the output file itself, and files that don't match extensions
		if d.IsDir() || filePath == outputPath || !shouldPrintFile(filePath, extensions) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			fmt.Printf("Error getting file info for %s: %v\n", filePath, err)
			return nil
		}
		fmt.Printf("%s (%d bytes)\n", filePath, info.Size())

		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			return nil
		}

		// Consolidate write operations
		if _, err := fmt.Fprintf(writer, "=== %s ===\n%s\n\n", filePath, content); err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", path, err)
		os.Exit(1)
	}

	fmt.Printf("Content written to %s\n", outputPath)
}

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

func shouldPrintFile(filePath string, extensions map[string]bool) bool {
	if extensions == nil {
		return true
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	return extensions[ext]
}
