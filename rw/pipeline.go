//
// pipeline.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  8 2014.
//

package rw

import (
	"errors"
	"fmt"
	"io"
)

type Pipeline struct {
	lastReader io.Reader
}

func (p *Pipeline) Exec(w io.WriteCloser) error {
	done := make(chan int)

	go func() {
		_, err := io.Copy(w, p.lastReader)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Pipeline successfull!!!")
		}
		done <- 1
	}()
	<-done
	return nil
}

func NewPipeline(encoders []Encoder, r io.Reader, d *Data) (*Pipeline, error) {
	if len(encoders) == 0 {
		return nil, errors.New("pipeline: A pipeline should have at lease 1 operation")
	}
	p := &Pipeline{
		lastReader: r,
	}
	for _, e := range encoders {
		nextReader, writer := io.Pipe()

		go func(e Encoder, r io.Reader, w *io.PipeWriter, d *Data) {
			e.Encode(r, w, d)
			w.Close()
		}(e, p.lastReader, writer, d)
		p.lastReader = nextReader
	}
	return p, nil
}
