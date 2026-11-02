# Move Duplicates Script (`movedups`)

`movedups` finds files with duplicate contents **in a single directory** and moves the extra copies into a separate directory. It keeps one copy of each matching group in the source directory.

It is a small Go command-line utility, not a recursive duplicate finder or a backup tool.

> [!CAUTION]
> **This command moves files immediately.** There is no preview, dry-run, confirmation prompt, or undo function. On systems where renaming over an existing file is allowed, a file with the same name in the destination **can be overwritten**. Back up important data, use a dedicated empty destination, and try the example below before using it on real files.

## Installation

Download a binary for your platform from [GitHub Releases](https://github.com/arran4/move-dups-script/releases), if one is available, or install with Go:

```sh
go install github.com/arran4/move-dups-script/cmd/movedups@latest
```

Make sure your Go binary installation directory (usually `$(go env GOPATH)/bin`) is on your `PATH`.

To build from a local checkout instead:

```sh
go build -o movedups ./cmd/movedups
```

## Quick start

To scan the files directly inside `./photos` and move extra copies to a separate sibling directory:

```sh
movedups -src ./photos -dest ./photos-duplicates
```

For example, suppose `photos` contains the following files:

```text
photos/
├── a.jpg     # contents A
├── b.jpg     # contents B
└── copy.jpg  # same contents as a.jpg
```

After the command completes successfully:

```text
photos/
├── a.jpg
└── b.jpg
photos-duplicates/
└── copy.jpg
```

File names do not need to match: duplicates are determined by **file size and MD5 digest of contents**, not by name. Within the source directory, the alphabetically first file in a matching group is kept.

### Try it on disposable files first

This example creates a temporary directory and does not touch your existing files (shell commands for Linux/macOS):

```sh
demo_dir="$(mktemp -d)"
mkdir -p "$demo_dir/source" "$demo_dir/duplicates"
printf 'same content\n' > "$demo_dir/source/a.txt"
cp "$demo_dir/source/a.txt" "$demo_dir/source/b.txt"
printf 'different content\n' > "$demo_dir/source/other.txt"

movedups -src "$demo_dir/source" -dest "$demo_dir/duplicates"
find "$demo_dir" -type f -print
```

You should find `a.txt` and `other.txt` in `source`, and `b.txt` in `duplicates`. Inspect `$demo_dir` and remove it when you are finished.

## Command-line options

```text
movedups [-src SOURCE_DIRECTORY] [-dest DESTINATION_DIRECTORY]
```

| Flag | Default | Meaning |
| --- | --- | --- |
| `-src` | `.` | Directory whose immediate files are checked for duplicates. |
| `-dest` | `dups` | Directory where extra copies are moved. Created if it does not exist. |
| `-h` | — | Show command-line help. |

With no arguments, `movedups` scans the **current working directory** and moves extra copies into `./dups`.

**Path resolution matters:** relative paths for both flags are resolved against the *working directory where the command is run*. In particular, `-dest dups` does **not** mean `<source>/dups` when you supply a different `-src`.

Use an explicit destination path to avoid surprises:

```sh
movedups -src "/path/to/photos" -dest "/path/to/photos-duplicates"
```

For safety, choose an initially empty destination outside the source directory, preferably on the same filesystem.

## How duplicate detection works

1. Read the direct children of the source directory, sorted by filename. **Subdirectories are not traversed.**
2. Get file sizes, skipping directory entries and logging entries whose metadata cannot be read.
3. Only hash files when at least two source entries have the same reported size. Files with unique sizes are left alone.
4. Read each candidate file and compute an MD5 digest. Group candidates by the number of bytes read and their digest.
5. Keep the first successfully hashed file in each group. Move later matches into the destination using their original filenames.

The destination directory is not scanned to discover duplicates. Different files with the same name are not duplicates unless their contents match; identical files with different names can be duplicates.

## Important limitations and precautions

- **Not recursive:** only files immediately inside `-src` are considered.
- **Not a dry run:** the tool moves files as it discovers duplicates. No review or rollback is provided.
- **Possible destination overwrite:** moves use the operating system's rename operation (`os.Rename`), which can replace an existing destination file on some systems. The program does not check for filename collisions or preserve pre-existing destination files.
- **Same-filesystem moves:** `os.Rename` can fail when source and destination are on different filesystems or volumes; there is no copy-and-delete fallback.
- **Partial failures may still exit successfully:** errors while reading, hashing, or moving individual files are logged and processing continues. Check the output and inspect both directories; a zero exit code does not guarantee every intended move succeeded.
- **MD5 is not a byte-for-byte verification:** a matching size and MD5 digest is treated as identical. MD5 collisions are possible, so do not use this as a high-assurance integrity check.
- **Avoid symbolic links and special files:** hashing follows a symbolic link to its target, but moving renames the link itself. Relative links may break after moving. This tool is best used on ordinary, unchanged files.
- **Do not modify files during a scan:** files changing between the size check, hashing, and move can produce unexpected results. There is no final content verification.

If you need recursive scanning, collision-safe quarantine, a dry-run, or stronger content verification, this version does not provide those features.

## Troubleshooting

**No files moved?** Check that duplicates are directly inside the source directory, that they have identical contents, and that you supplied the intended paths. Files with unique sizes are not hashed.

**A move failed?** Check permissions, destination filename conflicts, and whether the destination is on the same filesystem. Inspect the log messages for the affected file; other files may still have been processed.

**Where did the files go?** Check the `-dest` argument and your working directory. With default options, the destination is `./dups` relative to where the command was invoked.

## Development

Requires Go 1.22 or later.

```sh
go test ./...
```

The CLI entry point is in [`cmd/movedups`](cmd/movedups/main.go), and the duplicate-detection implementation is in [`pkg/movedups`](pkg/movedups/movedups.go).

## License

See [LICENSE](LICENSE) (BSD 3-Clause).
