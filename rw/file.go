//
// file.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  9 2014.
//

package rw

import (
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
		return RwError(s, "attribute 'dir' is empty")
	}
	if err := os.MkdirAll(s.Dir, 0755); err != nil {
		return RwError(s, err.Error())
	}
	return nil
}

func (s *File) NewWriter(id string, d *Data) (io.WriteCloser, error) {
	file, err := os.OpenFile(s.join(id), os.O_RDWR|os.O_CREATE, 0655)
	if err != nil {
		return nil, RwError(s, err.Error())
	}
	return file, nil
}

func (s *File) NewReader(id string, d *Data) (io.ReadCloser, error) {
	file, err := os.OpenFile(s.join(id), os.O_RDONLY, 0444)
	if err != nil {
		return nil, RwError(s, err.Error())
	}
	return file, nil
}

func (s *File) Delete(id string, d *Data) error {
	if err := os.Remove(s.join(id)); err != nil {
		return RwError(s, err.Error())
	}
	return nil
}
