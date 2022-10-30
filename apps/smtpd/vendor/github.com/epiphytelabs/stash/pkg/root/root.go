package root

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type FS string

func (fs FS) Create(path string) (*os.File, error) {
	file := fs.file(path)

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return nil, err
	}

	return os.Create(file)
}

func (fs FS) Exists(path string) (bool, error) {
	_, err := os.Stat(fs.file(path))
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (fs FS) Open(path string) (*os.File, error) {
	return os.Open(fs.file(path))
}

func (fs FS) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(fs.file(path))
}

func (fs FS) Remove(path string) error {
	return os.Remove(fs.file(path))
}

func (fs FS) Stat(path string) (os.FileInfo, error) {
	return os.Stat(fs.file(path))
}

func (fs FS) Write(path string, data []byte) error {
	file := fs.file(path)

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0644)
}

func (fs FS) file(path string) string {
	return filepath.Join(string(fs), path)
}
