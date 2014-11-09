//
// shell.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  5 2014.
//

package rw

import (
	"io"
	"os/exec"
)

type Shell struct {
	Cmd  string `json:"cmd"`
	Name string `json:"-"`
}

func (s *Shell) GetName() string {
	return s.Name
}

func (s *Shell) Init() error {
	_, err := exec.LookPath("bash")
	if err != nil {
		return RwError(s, "bash cannot be found in  current $PATH")
	}
	return nil
}

func (s *Shell) Encode(r io.Reader, w io.Writer, d *Data) error {
	cmd := exec.Command("bash", "-lc", s.Cmd)
	cmd.Stdout = w
	cmd.Stdin = r
	if err := cmd.Start(); err != nil {
		return RwError(s, err.Error())
	}

	if err := cmd.Wait(); err != nil {
		return RwError(s, err.Error())
	}
	return nil
}
