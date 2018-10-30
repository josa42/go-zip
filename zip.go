package zip

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type ProgressFunc func(path string, sourcePath string) bool

// Archive :
type Archive struct {
	filePath       string
	file           *os.File
	writer         *zip.Writer
	extractHandler func()
}

// CreateArchive :
func CreateArchive(path string) (Archive, error) {
	a := Archive{}

	a.filePath = path
	if err := a.create(path); err != nil {
		return a, err
	}

	return a, nil
}

// OpenArchive :
func OpenArchive(path string) (Archive, error) {
	a := Archive{}

	a.filePath = path
	if err := a.open(path); err != nil {
		return a, err
	}

	return a, nil
}

// Close :
func (a *Archive) Close() {

	if a.writer != nil {
		a.writer.Close()
	}

	if a.file != nil {
		a.file.Close()
	}
}

// Add :
func (a *Archive) Add(path, sourcePath string, progress ...ProgressFunc) error {
	if !a.isOpen() {
		return errors.New("archive is not open")
	}

	if path == "." || path == "" {
		path = "/"
	}

	progressFunc := func(path string, sourcePath string) bool { return true }
	if len(progress) == 1 {
		progressFunc = progress[0]
	}

	if !progressFunc(removeLeadingSlash(path), sourcePath) {
		return nil
	}

	if isDirectory(sourcePath) {
		return a.addDir(path, sourcePath, progressFunc)
	}

	return a.addFile(path, sourcePath)
}

// List :
func (a *Archive) List() ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(a.filePath)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()
		filenames = append(filenames, removeLeadingSlash(f.Name))
	}

	return filenames, nil
}

func (a *Archive) isOpen() bool {
	return a.file != nil && a.writer != nil
}

func (a *Archive) create(path string) error {
	if a.isOpen() {
		return errors.New("archive is already open")
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	a.file = file
	a.writer = zip.NewWriter(a.file)

	return nil
}

func (a *Archive) open(path string) error {
	if a.isOpen() {
		return errors.New("archive is already open")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	a.file = file
	a.writer = zip.NewWriter(a.file)

	return nil
}

func (a *Archive) addDir(path, sourcePath string, progress ProgressFunc) error {

	fileInfos, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	if path != "/" {
		a.addFile(path+"/", sourcePath)
	}
	for _, fileInfo := range fileInfos {
		a.Add(filepath.Join(path, fileInfo.Name()), filepath.Join(sourcePath, fileInfo.Name()), progress)
	}

	return err
}

func (a *Archive) addFile(path, sourcePath string) error {
	if !a.isOpen() {
		return errors.New("archive is not open")
	}

	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if path == "." {
		path = filepath.Join(path, filepath.Base(sourcePath))
	}

	fileWriter, err := a.writer.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, file)

	a.writer.Flush()

	return err
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func removeLeadingSlash(path string) string {
	var re = regexp.MustCompile(`^/`)
	return re.ReplaceAllString(path, "")
}
