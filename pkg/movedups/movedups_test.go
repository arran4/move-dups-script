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
