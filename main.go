package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func main() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Panic(err)
	}
	if err := os.Mkdir("dups", 0755); err != nil {
		log.Printf("Error creating %s because %v", "dups", err)
	}
	seenHashes := map[string]struct{}{}
	for fileN, file := range files {
		fileAction(fileN, file, seenHashes, len(files))
	}
}

func fileAction(fileN int, file os.FileInfo, seenHashes map[string]struct{}, fileCount int) {
	log.Printf("File %d of %d", fileN, fileCount)
	if file.IsDir() {
		return
	}
	mh := md5.New()
	fh, err := os.Open(file.Name())
	if err != nil {
		log.Printf("Error with %s: %v", file.Name(), err)
		return
	}
	defer fh.Close()
	n, err := io.Copy(mh, fh)
	if err != nil {
		log.Printf("Error with %s: %v", file.Name(), err)
		return
	}
	if n != file.Size() {
		log.Printf("Note: %s says it's %d but read %d differece of %d", file.Name(), n, file.Size(), n-file.Size())
	}
	sum := hex.EncodeToString(mh.Sum(nil))
	hash := fmt.Sprintf("%d-%s", n, sum)
	if _, ok := seenHashes[hash]; ok {
		newPath := path.Join("dups", file.Name())
		log.Printf("Moving %s to %s", file.Name(), newPath)
		if err := os.Rename(file.Name(), newPath); err != nil {
			log.Printf("Error moving file: %s because %v", file.Name(), err)
		}
	}
	seenHashes[hash] = struct{}{}
}
