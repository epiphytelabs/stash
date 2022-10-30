package root

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type FS string

func (fs FS) Create(path string) (*os.File, error) {
	file := fs.file(path)

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return nil, errors.WithStack(err)
	}

	fd, err := os.Create(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return fd, nil
}

func (fs FS) Exists(path string) (bool, error) {
	_, err := os.Stat(fs.file(path))
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, errors.WithStack(err)
	}

	return true, nil
}

func (fs FS) Open(path string) (*os.File, error) {
	fd, err := os.Open(fs.file(path))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return fd, nil
}

func (fs FS) Read(path string) ([]byte, error) {
	data, err := os.ReadFile(fs.file(path))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func (fs FS) Remove(path string) error {
	if err := os.Remove(fs.file(path)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (fs FS) Stat(path string) (os.FileInfo, error) {
	info, err := os.Stat(fs.file(path))
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return info, nil
}

func (fs FS) Write(path string, data []byte) error {
	file := fs.file(path)

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return errors.WithStack(err)
	}

	if err := os.WriteFile(file, data, 0644); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (fs FS) file(path string) string {
	return filepath.Join(string(fs), path)
}
