/*
Copyright (c) 2025, arran4
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package movedups

import (
	"crypto/md5"
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

	seenHashes := map[fileKey]struct{}{}
	fileCount := len(entries)

	fileSizes := make([]int64, len(entries))
	sizeCounts := make(map[int64]int)

	for i, entry := range entries {
		fileSizes[i] = -1
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			log.Printf("Error getting info for %s: %v", filepath.Join(srcDir, entry.Name()), err)
			continue
		}
		size := info.Size()
		fileSizes[i] = size
		sizeCounts[size]++
	}

	for i, entry := range entries {
		size := fileSizes[i]
		if size == -1 {
			continue
		}
		if sizeCounts[size] <= 1 {
			continue
		}
		processFile(i, entry, srcDir, destDir, seenHashes, fileCount, size)
	}
	return nil
}

type fileKey struct {
	size int64
	hash [16]byte
}

func processFile(fileN int, entry os.DirEntry, srcDir, destDir string, seenHashes map[fileKey]struct{}, fileCount int, size int64) {
	log.Printf("File %d of %d: %s", fileN+1, fileCount, entry.Name())
	if entry.IsDir() {
		return
	}

	fullPath := filepath.Join(srcDir, entry.Name())

	hash, err := calculateHash(fullPath, size)
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

func calculateHash(path string, expectedSize int64) (fileKey, error) {
	fh, err := os.Open(path)
	if err != nil {
		return fileKey{}, err
	}
	defer fh.Close()

	mh := md5.New()
	n, err := io.Copy(mh, fh)
	if err != nil {
		return fileKey{}, err
	}

	if n != expectedSize {
		log.Printf("Note: %s says it's %d but read %d difference of %d", path, n, expectedSize, n-expectedSize)
	}

	var k fileKey
	k.size = n
	mh.Sum(k.hash[:0])
	return k, nil
}
