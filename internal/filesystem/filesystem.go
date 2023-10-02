package filesystem

import (
	"os"
	"strings"
)

// FileWriter describes an interface for writing files to the file system.
type FileWriter interface {
	WriteFile(path string, data []byte) error
}

// fileWriter implements the FileWriter interface.
type fileWriter struct{}

func NewFileWriter() fileWriter {
	return fileWriter{}
}

// WriteFile writes a file to the filesystem, if the directories in the path don't exist, they will be created.
func (fw fileWriter) WriteFile(path string, data []byte) error {
	if err := os.MkdirAll(getDirectory(path), 0644); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func getDirectory(path string) string {
	l := strings.LastIndex(path, "/")

	return path[:l]
}
