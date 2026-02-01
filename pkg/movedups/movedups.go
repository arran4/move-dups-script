package movedups

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type job struct {
	index int
	entry os.DirEntry
}

type result struct {
	index int
	entry os.DirEntry
	hash  string
	err   error
	skip  bool
}

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

	// Start workers
	numWorkers := runtime.NumCPU()

	jobs := make(chan job, numWorkers*2)
	results := make(chan result, numWorkers*2)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(srcDir, jobs, results)
		}()
	}

	// Send jobs
	go func() {
		for i, entry := range entries {
			jobs <- job{index: i, entry: entry}
		}
		close(jobs)
	}()

	// Wait for workers to finish and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Process results in order
	nextIndex := 0
	resultBuffer := make(map[int]result)

	for res := range results {
		resultBuffer[res.index] = res

		for {
			if r, ok := resultBuffer[nextIndex]; ok {
				delete(resultBuffer, nextIndex)
				processResult(r, srcDir, destDir, seenHashes, fileCount)
				nextIndex++
			} else {
				break
			}
		}
	}
	return nil
}

func worker(srcDir string, jobs <-chan job, results chan<- result) {
	for j := range jobs {
		res := result{index: j.index, entry: j.entry}

		if j.entry.IsDir() {
			res.skip = true
			results <- res
			continue
		}

		fullPath := filepath.Join(srcDir, j.entry.Name())

		// Get size for hash key and verification
		info, err := j.entry.Info()
		if err != nil {
			res.err = fmt.Errorf("getting info for %s: %w", fullPath, err)
			results <- res
			continue
		}

		hash, err := calculateHash(fullPath, info.Size())
		if err != nil {
			res.err = fmt.Errorf("hashing %s: %w", fullPath, err)
			results <- res
			continue
		}

		res.hash = hash
		results <- res
	}
}

func processResult(r result, srcDir, destDir string, seenHashes map[string]struct{}, fileCount int) {
	log.Printf("File %d of %d: %s", r.index+1, fileCount, r.entry.Name())

	if r.skip {
		return
	}

	if r.err != nil {
		log.Printf("Error: %v", r.err)
		return
	}

	fullPath := filepath.Join(srcDir, r.entry.Name())

	if _, ok := seenHashes[r.hash]; ok {
		newPath := filepath.Join(destDir, r.entry.Name())
		log.Printf("Moving %s to %s", fullPath, newPath)
		if err := os.Rename(fullPath, newPath); err != nil {
			log.Printf("Error moving file: %s because %v", fullPath, err)
		}
	} else {
		seenHashes[r.hash] = struct{}{}
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
