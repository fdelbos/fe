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

type EncodingPipeline struct {
	Encoders []Encoder
	Output   Outputer
}

type DecodingPipeline struct {
	Decoders []Decoder
	Input    Inputer
}

func (p *Pipeline) Exec(w io.Writer) error {
	done := make(chan int)

	go func() {
		_, err := io.Copy(w, p.lastReader)

		if err != nil {
			fmt.Println(err)
		}
		done <- 1
	}()
	<-done

	return nil
}

func NewEncoding(encoders []Encoder, r io.Reader, d *Data) (*Pipeline, error) {
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

func NewDecoding(decoders []Decoder, r io.Reader, d *Data) (*Pipeline, error) {
	if len(decoders) == 0 {
		return nil, errors.New("pipeline: A pipeline should have at lease 1 operation")
	}
	p := &Pipeline{
		lastReader: r,
	}
	for _, e := range decoders {
		nextReader, writer := io.Pipe()

		go func(e Decoder, r io.Reader, w *io.PipeWriter, d *Data) {
			e.Decode(r, w, d)
			w.Close()
		}(e, p.lastReader, writer, d)

		p.lastReader = nextReader
	}
	return p, nil
}
