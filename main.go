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
	flag.Usage = func() {
		fmt.Println("👏 Clap slaps all your files into one!")
		fmt.Println("Usage: clap [-o filename] <path> [ext <extensions...>] [sub <substrings...>]")
		fmt.Println("  -o filename  output filename (default: clap.txt in current working directory)")
		fmt.Println("Filters are case-insensitive. Extensions may include or omit the leading dot.")
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		fmt.Println("Error: path is required")
		os.Exit(1)
	}

	path := args[0]
	extFilters, subFilters, err := parseFilterArgs(args[1:])
	if err != nil {
		fmt.Printf("Error parsing filters: %v\n", err)
		os.Exit(1)
	}
	outputPath, err := resolveOutputPath(*outputFilename)
	if err != nil {
		fmt.Printf("Error resolving output path: %v\n", err)
		os.Exit(1)
	}

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

		// Skip directories, the output file itself, and files that don't match filters
		if d.IsDir() {
			return nil
		}

		fileAbs, absErr := filepath.Abs(filePath)
		if absErr == nil && fileAbs == outputPath {
			return nil
		}

		if !shouldPrintFile(filePath, extFilters, subFilters) {
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

func resolveOutputPath(outputFilename string) (string, error) {
	if filepath.IsAbs(outputFilename) {
		return filepath.Clean(outputFilename), nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Clean(filepath.Join(cwd, outputFilename)), nil
}

func parseFilterArgs(args []string) ([]string, []string, error) {
	var extFilters []string
	var subFilters []string

	var state string
	var lastKeyword string
	var lastKeywordHasValue bool

	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}

		switch arg {
		case "ext", "sub":
			if lastKeyword != "" && !lastKeywordHasValue {
				return nil, nil, fmt.Errorf("%s requires at least one value", lastKeyword)
			}
			state = arg
			lastKeyword = arg
			lastKeywordHasValue = false
		default:
			if state == "" {
				return nil, nil, fmt.Errorf("unexpected filter value %q; expected \"ext\" or \"sub\"", arg)
			}
			if state == "ext" {
				normalized := normalizeExt(arg)
				if normalized == "" {
					continue
				}
				extFilters = append(extFilters, normalized)
			} else {
				subFilters = append(subFilters, strings.ToLower(arg))
			}
			lastKeywordHasValue = true
		}
	}

	if lastKeyword != "" && !lastKeywordHasValue {
		return nil, nil, fmt.Errorf("%s requires at least one value", lastKeyword)
	}

	return extFilters, subFilters, nil
}

func normalizeExt(ext string) string {
	ext = strings.TrimSpace(ext)
	ext = strings.TrimPrefix(ext, ".")
	if ext == "" {
		return ""
	}
	return strings.ToLower(ext)
}

func shouldPrintFile(filePath string, extFilters, subFilters []string) bool {
	if len(extFilters) == 0 && len(subFilters) == 0 {
		return true
	}

	fileName := filepath.Base(filePath)
	lowerName := strings.ToLower(fileName)

	extOK := len(extFilters) == 0
	if !extOK {
		fileExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")
		for _, filter := range extFilters {
			if filter == fileExt {
				extOK = true
				break
			}
		}
	}

	subOK := len(subFilters) == 0
	if !subOK {
		for _, filter := range subFilters {
			if strings.Contains(lowerName, filter) {
				subOK = true
				break
			}
		}
	}

	return extOK && subOK
}
