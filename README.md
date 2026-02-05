# Clap

CLI tool that combines multiple files from a directory into a single output file.

## Build

Requirements: Go.

```bash
git clone https://github.com/yourusername/clap.git
cd clap
go build
```

## Usage

```bash
clap [-o filename] <path> [filters...]
```

- `-o filename` sets the output file name (default: `clap.txt`). The file is created in `<path>`.
- `filters` are case-sensitive substrings matched against each file’s base name (not just extensions).

### Examples

Full snapshot (all files):

```bash
clap ./myproject
```

Code review bundle (source files only):

```bash
clap -o review.txt ./myproject .go .md
```

LLM context file (source + docs):

```bash
clap -o context.txt ./src .js .jsx .ts .tsx .md
```

Documentation bundle:

```bash
clap -o all-docs.md ./docs .md
```

Targeted analysis by filename substring:

```bash
clap ./myproject _test .config
```

### Output Format

```
=== path/to/file1.go ===
[file content]

=== path/to/file2.go ===
[file content]
```
