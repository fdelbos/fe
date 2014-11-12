package main

import (
	"bytes"
	"crypto/rand"
	"github.com/dchest/uniuri"
	"github.com/fdelbos/fe/rw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("Service", func() {

	testBin := make([]byte, 1<<16)
	rand.Read(testBin)
	data := rw.NewData()
	id := uniuri.New()
	cat := &rw.Shell{
		Cmd:  "cat",
		Name: "cat",
	}
	cat.Init()

	It("run a simple pipeline", func() {
		out := new(bytes.Buffer)
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p := rw.NewEncoding(
			[]rw.Encoder{cat},
			r,
			data)
		Ω(p.Exec(out)).To(BeNil())
		Ω(bytes.Equal(out.Bytes(), testBin)).To(BeTrue())
	})

	file := &rw.File{
		Dir:  "/tmp/" + uniuri.New(),
		Name: "file",
	}
	file.Init()

	service := &Service{
		Url: "/test",
		EncodingPipe: &rw.EncodingPipeline{
			Output: file,
		},
	}

	It("should encode", func() {
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		Ω(service.Encode(id, r, data)).To(BeNil())
	})

	service.DecodingPipe = &rw.DecodingPipeline{
		Input: file,
	}
	It("should decode", func() {
		out := new(bytes.Buffer)
		r, w := io.Pipe()
		go func() {
			io.Copy(out, r)
		}()
		Ω(service.Decode(id, w, data)).To(BeNil())
		Ω(bytes.Equal(out.Bytes(), testBin)).To(BeTrue())
	})

})
