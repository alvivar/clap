// Package declaration: Every Go file starts with a package name.
// "main" is special - it tells Go this is an executable program, not a library.
// The "main" package must have a "main()" function as the entry point.
package main

// Import section: Brings in standard library packages we need.
// Go uses double quotes for imports, and each import is on its own line in parentheses.
import (
	"flag"          // Command-line flag parsing (like -o filename)
	"fmt"           // Formatted I/O (printing to console)
	"os"            // Operating system functions (file I/O, exit codes)
	"path/filepath" // Cross-platform file path manipulation
	"strings"       // String manipulation utilities
)

// main is the entry point of the program. Every executable Go program needs one.
// It takes no parameters and returns nothing (uses os.Exit for exit codes).
func main() {
	// Variable declaration with :=
	// The := operator declares and initializes in one step (type is inferred).
	// flag.String() returns a *string (pointer to string).
	// This defines a command-line flag: -o <filename>
	outputFilename := flag.String("o", "clap.file", "output filename")

	// Parse the command-line flags. Must be called before accessing flag values.
	flag.Parse()

	// flag.Args() returns remaining arguments after flags are parsed.
	// Here, we expect: <path> [extensions...]
	args := flag.Args()

	// len() is a built-in function that returns the length of slices, arrays, strings, etc.
	// Go uses simple if statements without parentheses around the condition.
	// Check that at least the path argument is provided (path is required)
	if len(args) < 1 {
		// fmt.Println prints to stdout with a newline
		fmt.Println("👏 Clap slaps all your files into one!")
		fmt.Println("Usage: clap [-o filename] <path> [extensions...]")
		fmt.Println("Error: path is required")
		// os.Exit terminates the program with an exit code (1 = error)
		os.Exit(1)
	}

	// Slice indexing: args[0] gets the first element
	path := args[0]
	// Slice slicing: args[1:] creates a new slice from index 1 to the end
	extensions := normalizeExtensions(args[1:])

	// Variable declaration with var keyword and type inference.
	// strings.Builder is efficient for building strings incrementally (avoids allocations).
	var contentBuilder strings.Builder

	// filepath.Walk recursively walks a directory tree.
	// It takes a path and a callback function (anonymous function/closure here).
	// The callback receives: path, file info, and any error encountered.
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		// Always check errors! Go uses explicit error handling (no exceptions).
		// If there was an error accessing this path, print it and return the error.
		if err != nil {
			// fmt.Printf is like printf in C - %s for string, %v for any value
			fmt.Printf("Error accessing path %s: %v\n", filePath, err)
			return err
		}

		// Skip directories and files that don't match our extension filter.
		// The || operator means "or" (Go uses || for or, && for and, ! for not).
		if info.IsDir() || !shouldPrintFile(filePath, extensions) {
			// Returning nil means "no error, continue walking"
			return nil
		}

		// Print the file we're processing (%d is for integers)
		fmt.Printf("%s (%d bytes)\n", filePath, info.Size())

		// os.ReadFile reads the entire file into a byte slice.
		// Multiple return values: Go functions can return multiple values.
		// By convention, the last return value is often an error.
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			// Return nil to continue processing other files despite this error
			return nil
		}

		// Build the output format: === filepath ===\ncontents\n\n
		// WriteString appends a string to the builder
		contentBuilder.WriteString("=== ")
		contentBuilder.WriteString(filePath)
		contentBuilder.WriteString(" ===\n")
		// Write appends a byte slice (content is []byte from ReadFile)
		contentBuilder.Write(content)
		contentBuilder.WriteString("\n\n")

		// Return nil to indicate success and continue walking
		return nil
	})

	// Check if the Walk operation itself failed
	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", path, err)
		os.Exit(1)
	}

	// filepath.Join concatenates path elements with the OS-specific separator.
	// The * dereferences the pointer to get the actual string value.
	outputPath := filepath.Join(path, *outputFilename)

	// os.WriteFile writes data to a file with specified permissions.
	// []byte() converts the string to a byte slice.
	// 0644 is Unix-style permissions (owner: rw, group: r, others: r).
	// The if statement can include initialization: if err := ...; err != nil
	if err := os.WriteFile(outputPath, []byte(contentBuilder.String()), 0644); err != nil {
		fmt.Printf("Error writing output file %s: %v\n", outputPath, err)
		os.Exit(1)
	}

	fmt.Printf("Content written to %s\n", outputPath)
}

// Function declaration: func keyword, name, parameters in parentheses, return type.
// This function takes a slice of strings and returns a map[string]bool (or nil).
// normalizeExtensions converts extensions to a map with leading dots and lowercase.
// Returns nil if no extensions provided (accept all files).
func normalizeExtensions(extensions []string) map[string]bool {
	// Early return pattern: handle edge cases first
	if len(extensions) == 0 {
		// nil is the zero value for pointers, slices, maps, channels, etc.
		return nil
	}

	// make() is a built-in function that creates slices, maps, and channels.
	// Here we create a map with string keys and bool values.
	// The second argument is the capacity hint (optional optimization).
	// Maps in Go are hash tables - O(1) average lookup time.
	extMap := make(map[string]bool, len(extensions))

	// for...range loop iterates over slices, arrays, maps, strings, etc.
	// It returns two values: index and value.
	// The _ (blank identifier) discards the index since we don't need it.
	for _, ext := range extensions {
		// strings.HasPrefix checks if a string starts with a prefix
		if !strings.HasPrefix(ext, ".") {
			// String concatenation with + operator
			// Note: ext is reassigned here (shadows the loop variable in this scope)
			ext = "." + ext
		}
		// Map assignment: set the value for this key
		// Using bool map as a "set" data structure (only care about keys)
		extMap[strings.ToLower(ext)] = true
	}

	// Named return: the function signature declares what we return
	return extMap
}

// shouldPrintFile returns true if the file matches the extension filter.
// If extensions is nil, all files are included.
// This function demonstrates Go's simple return type syntax: just "bool" at the end.
func shouldPrintFile(filePath string, extensions map[string]bool) bool {
	// Checking if a map is nil (no filter specified means accept all files)
	if extensions == nil {
		return true
	}

	// filepath.Ext extracts the file extension including the dot (e.g., ".go")
	// strings.ToLower converts to lowercase for case-insensitive matching
	ext := strings.ToLower(filepath.Ext(filePath))

	// Map lookup: extensions[ext] returns the value and a boolean "ok" flag.
	// Here we only use the value. If the key doesn't exist, Go returns the
	// zero value for the type (false for bool), which is perfect for our use case.
	// This is effectively checking "is ext in the set?"
	return extensions[ext]
}
