package movedups

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	// Create source directory
	srcDir := t.TempDir()

	// Create destination directory
	destDir := t.TempDir()

	// Create file 1
	file1Content := []byte("hello world")
	if err := os.WriteFile(filepath.Join(srcDir, "file1.txt"), file1Content, 0644); err != nil {
		t.Fatal(err)
	}

	// Create file 2 (different content)
	file2Content := []byte("hello universe")
	if err := os.WriteFile(filepath.Join(srcDir, "file2.txt"), file2Content, 0644); err != nil {
		t.Fatal(err)
	}

	// Create file 3 (duplicate of file 1)
	if err := os.WriteFile(filepath.Join(srcDir, "file3.txt"), file1Content, 0644); err != nil {
		t.Fatal(err)
	}

	// Create file 4 (duplicate of file 1)
	if err := os.WriteFile(filepath.Join(srcDir, "file4.txt"), file1Content, 0644); err != nil {
		t.Fatal(err)
	}

	// Run logic
	if err := Run(srcDir, destDir); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify srcDir contains file1.txt and file2.txt, but NOT file3.txt and file4.txt
	// Iteration order of `os.ReadDir` is by name.
	// So file1.txt (first) -> kept.
	// file2.txt -> kept.
	// file3.txt -> duplicate of file1 -> moved.
	// file4.txt -> duplicate of file1 -> moved.

	assertFileExists(t, filepath.Join(srcDir, "file1.txt"))
	assertFileExists(t, filepath.Join(srcDir, "file2.txt"))
	assertFileMissing(t, filepath.Join(srcDir, "file3.txt"))
	assertFileMissing(t, filepath.Join(srcDir, "file4.txt"))

	// Verify destDir contains file3.txt and file4.txt
	assertFileExists(t, filepath.Join(destDir, "file3.txt"))
	assertFileExists(t, filepath.Join(destDir, "file4.txt"))
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s", path)
	}
}

func assertFileMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Errorf("Expected file to be missing: %s", path)
	}
}
