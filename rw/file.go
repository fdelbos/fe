//
// file.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  9 2014.
//

package rw

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	Dir  string `json:"dir"`
	Name string `json:"name"`
}

func (s *File) join(path string) string {
	return filepath.Join(s.Dir, path)
}

func (s *File) GetName() string {
	return s.Name
}

func (s *File) Init() error {
	if s.Dir == "" {
		return errors.New("attribute 'dir' is empty")
	}
	if err := os.MkdirAll(s.Dir, 0755); err != nil {
		return err
	}
	return nil
}

func (s *File) NewWriter(id string, d *Data) (io.WriteCloser, error) {
	return os.OpenFile(s.join(id), os.O_RDWR|os.O_CREATE, 0655)
}

func (s *File) NewReader(id string, d *Data) (io.ReadCloser, error) {
	return os.OpenFile(s.join(id), os.O_RDONLY, 0444)
}

func (s *File) Delete(id string, d *Data) error {
	return os.Remove(s.join(id))
}
