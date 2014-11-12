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

func (s *File) NewWriter(d *Data) (io.WriteCloser, error) {
	id := d.Get("identifier").(string)
	if id == "" {
		return nil, RwError(s, "No identifier")
	}
	file, err := os.OpenFile(s.join(id), os.O_RDWR|os.O_CREATE, 0655)
	if err != nil {
		return nil, RwError(s, err.Error())
	}
	return file, nil
}

func (s *File) NewReader(d *Data) (io.ReadCloser, error) {
	id := d.Get("identifier").(string)
	if id == "" {
		return nil, RwError(s, "No identifier")
	}
	file, err := os.OpenFile(s.join(id), os.O_RDONLY, 0444)
	if err != nil {
		return nil, RwError(s, err.Error())
	}
	return file, nil
}

func (s *File) Delete(d *Data) error {
	id := d.Get("identifier").(string)
	if id == "" {
		return RwError(s, "No identifier")
	}
	if err := os.Remove(s.join(id)); err != nil {
		return RwError(s, err.Error())
	}
	return nil
}
