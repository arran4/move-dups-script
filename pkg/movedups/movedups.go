package movedups

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Run(srcDir, destDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("reading source directory: %w", err)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("creating destination directory: %w", err)
	}

	seenHashes := map[string]struct{}{}
	fileCount := len(entries)

	for i, entry := range entries {
		processFile(i, entry, srcDir, destDir, seenHashes, fileCount)
	}
	return nil
}

func processFile(fileN int, entry os.DirEntry, srcDir, destDir string, seenHashes map[string]struct{}, fileCount int) {
	log.Printf("File %d of %d: %s", fileN+1, fileCount, entry.Name())
	if entry.IsDir() {
		return
	}

	fullPath := filepath.Join(srcDir, entry.Name())

	// Get size for hash key and verification
	info, err := entry.Info()
	if err != nil {
		log.Printf("Error getting info for %s: %v", fullPath, err)
		return
	}

	hash, err := calculateHash(fullPath, info.Size())
	if err != nil {
		log.Printf("Error hashing %s: %v", fullPath, err)
		return
	}

	if _, ok := seenHashes[hash]; ok {
		newPath := filepath.Join(destDir, entry.Name())
		log.Printf("Moving %s to %s", fullPath, newPath)
		if err := os.Rename(fullPath, newPath); err != nil {
			log.Printf("Error moving file: %s because %v", fullPath, err)
		}
	} else {
		seenHashes[hash] = struct{}{}
	}
}

func calculateHash(path string, expectedSize int64) (string, error) {
	fh, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	mh := md5.New()
	n, err := io.Copy(mh, fh)
	if err != nil {
		return "", err
	}

	if n != expectedSize {
		log.Printf("Note: %s says it's %d but read %d difference of %d", path, n, expectedSize, n-expectedSize)
	}

	sum := hex.EncodeToString(mh.Sum(nil))
	return fmt.Sprintf("%d-%s", n, sum), nil
}
