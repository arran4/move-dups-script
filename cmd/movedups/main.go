package main

import (
	"flag"
	"log"
	"os"

	"github.com/arran4/move-dups-script/pkg/movedups"
)

func main() {
	srcDir := flag.String("src", ".", "Source directory to scan for duplicates")
	destDir := flag.String("dest", "dups", "Destination directory to move duplicates to")
	flag.Parse()

	if *srcDir == "" || *destDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("Scanning %s and moving duplicates to %s", *srcDir, *destDir)

	if err := movedups.Run(*srcDir, *destDir); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
