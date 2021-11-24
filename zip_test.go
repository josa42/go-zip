package zip_test

import (
	"os"
	"testing"

	zip "github.com/josa42/go-zip"
)

func TestCreateArchive(t *testing.T) {
	// t.Skip("skipping test in short mode.")

	a, err := zip.CreateArchive("test.zip")
	a.Close()

	assertErrorIsNil(t, err)
	assertFileExists(t, "test.zip")

	removeFile("test.zip")
}

func TestArchiveAdd(t *testing.T) {

	a, err := zip.CreateArchive("test.zip")
	a.Add("/", "testdata/source2")
	a.Close()

	assertErrorIsNil(t, err)
	assertFileExists(t, "test.zip")
	assertArchiveHasContent(t, "test.zip", []string{"not-ignored2"})

	removeFile("test.zip")
}

func TestArchiveAddIngore(t *testing.T) {
	a, err := zip.CreateArchive("test.zip")
	a.Add("/", "testdata/source1", func(path string, sourcePath string) bool {
		return path != ".ignored"
	})
	a.Close()

	assertErrorIsNil(t, err)
	assertFileExists(t, "test.zip")
	assertArchiveHasContent(t, "test.zip", []string{"dir/", "dir/not-ignored", "not-ignored"})

	removeFile("test.zip")
}

func TestArchiveList(t *testing.T) {
	a, err1 := zip.OpenArchive("testdata/archive.zip")
	list, err2 := a.List()
	a.Close()

	assertErrorIsNil(t, err1)
	assertErrorIsNil(t, err2)
	assertListHasContent(t, list, []string{"not-ignored2"})

	removeFile("test.zip")
}

func assertFileExists(t *testing.T, filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("file should exist: %s", filePath)
	}
}
func assertErrorIsNil(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("err should be nil: %v", err)
	}
}

func assertArchiveHasContent(t *testing.T, path string, content []string) {
	a, _ := zip.OpenArchive(path)
	list, _ := a.List()

	assertListHasContent(t, list, content)
}

func assertListHasContent(t *testing.T, list, content []string) {
	if len(list) != len(content) {
		t.Fatalf("should contain %d items (not %d, %v)", len(content), len(list), list)
	}

	for idx, item := range content {
		if list[idx] != item {
			t.Fatalf("item at index %d should be %s (not %s)", idx, item, list[idx])
		}
	}
}

func removeFile(filePath string) {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		os.Remove(filePath)
	}
}
