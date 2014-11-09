//
// rw.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  8 2014.
//

package rw

import (
	"errors"
	"fmt"
	"io"
)

type DataMap map[string]interface{}

type RwBase interface {
	GetName() string
	Init() error
}

type Encoder interface {
	RwBase
	Encode(io.Reader, io.Writer, *Data) error
}

type Decoder interface {
	RwBase
	Decode(io.Reader, io.Writer, *Data) error
}

type EncodeDecoder interface {
	RwBase
	Encode(io.Reader, io.Writer, *Data) error
	Decode(io.Reader, io.Writer, *Data) error
}

type Outputer interface {
	RwBase
	NewWirter() io.Writer
}

type Inputer interface {
	RwBase
	NewReader() io.Reader
}

func RwError(e RwBase, err string) error {
	fmt.Println(fmt.Sprintf("%s: %s", e.GetName(), err))
	return errors.New(fmt.Sprintf("%s: %s", e.GetName(), err))
}
