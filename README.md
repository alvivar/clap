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
clap [-o filename] <path> [ext <extensions...>] [sub <substrings...>]
```

- `-o filename` sets the output file name (default: `clap.txt`). The file is created in the current working directory unless `-o` includes a path.
- `ext` filters by file extension (case-insensitive). The dot prefix is optional.
- `sub` filters by case-insensitive substring matches against each file’s base name.
- When both `ext` and `sub` are provided, files must match both.
- Paths written in output headers are relative to the scanned `<path>` root.

After `<path>`, you can include an `ext` section and/or a `sub` section. Each section accepts multiple space-separated values.

### Examples

Full snapshot (all files):

```bash
clap ./myproject
```

Code review bundle (source files only):

```bash
clap -o review.txt ./myproject ext go md
```

LLM context file (source + docs):

```bash
clap -o context.txt ./src ext js jsx ts tsx md
```

Documentation bundle:

```bash
clap -o all-docs.md ./docs ext md
```

Targeted analysis by filename substring:

```bash
clap ./myproject sub _test .config
```

Combined extension + substring filtering:

```bash
clap ./myproject ext go sub test
```

### Output Format

```
=== path/to/file1.go ===
[file content]

=== path/to/file2.go ===
[file content]
```
