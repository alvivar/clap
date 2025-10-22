# 👏 Clap

**Clap slaps all your files into one!**

A blazingly fast CLI tool that combines multiple files from a directory into a single, organized output file. Perfect for sharing codebases with LLMs, conducting code reviews, or creating comprehensive documentation snapshots.

## ✨ Features

-   🚀 **Fast & Efficient** - Recursively walks through directories at lightning speed
-   🎯 **Smart Filtering** - Filter files by extension (supports multiple extensions)
-   📝 **Clear Formatting** - Each file is clearly separated with headers showing the file path
-   💪 **Flexible Output** - Customize the output filename to your needs
-   🔄 **Case-Insensitive** - Extensions work with or without dots (`.go` or `go`)
-   📊 **Progress Tracking** - See which files are being processed with size information
-   🌳 **Recursive Search** - Automatically traverses nested directories

## 🚀 Installation

```bash
go install github.com/yourusername/clap@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/clap.git
cd clap
go build
```

## 📖 Usage

### Basic Usage

Combine all files in a directory:

```bash
clap /path/to/directory
```

### Filter by Extensions

Combine only specific file types:

```bash
# Single extension
clap /path/to/project .go

# Multiple extensions
clap /path/to/project .go .md .txt

# Extensions work without dots too!
clap /path/to/project go md txt
```

### Custom Output File

Specify a custom output filename:

```bash
clap -o combined.txt /path/to/directory .js .ts
```

## 📚 Examples

**Combine all Go files in a project:**

```bash
clap ./myproject .go
```

**Create a codebase snapshot for AI:**

```bash
clap -o context.txt ./src .js .jsx .ts .tsx
```

**Gather all documentation:**

```bash
clap -o all-docs.md ./docs .md
```

## 📋 Output Format

Clap creates a well-organized output file with clear separators:

```
=== path/to/file1.go ===
[file content]

=== path/to/file2.go ===
[file content]
```

## 🎯 Use Cases

-   **AI Context Building** - Feed entire codebases to Large Language Models
-   **Code Reviews** - Share complete project snapshots
-   **Documentation** - Combine multiple markdown files into one
-   **Backup** - Create text-based snapshots of your projects
-   **Code Analysis** - Prepare files for analysis tools

## 📝 License

MIT License - feel free to use this in your projects!
