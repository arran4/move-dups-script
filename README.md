# Move Duplicates Script

A simple Go tool to scan a directory for duplicate files (based on size and MD5 hash) and move them to a separate directory.

## Installation

You can download the latest release from the [Releases](https://github.com/arran4/move-dups-script/releases) page.

Or install with Go:

```bash
go install github.com/arran4/move-dups-script/cmd/movedups@latest
```

## Usage

```bash
movedups -src <source_directory> -dest <destination_directory>
```

### Flags

- `-src`: Source directory to scan for duplicates. Defaults to `.` (current directory).
- `-dest`: Destination directory to move duplicates to. Defaults to `dups`.

### Examples

Scan the current directory and move duplicates to `dups`:

```bash
movedups
```

Scan the `photos` directory and move duplicates to `duplicates`:

```bash
movedups -src ./photos -dest ./duplicates
```

## How it works

1. Scans the source directory for files.
2. Calculates the MD5 hash of each file (including size in the hash key).
3. Keeps track of seen hashes.
4. If a file's hash has been seen before, it moves the file to the destination directory.
5. If the destination directory does not exist, it is created.

## Development

### Prerequisites

- Go 1.22 or later

### Running Tests

```bash
go test ./...
```
