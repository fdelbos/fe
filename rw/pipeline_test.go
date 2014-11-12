package rw_test

import (
	"bytes"
	"crypto/rand"
	"io"
	. "github.com/fdelbos/fe/rw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipeline", func() {

	testBin := make([]byte, 1<<18)
	rand.Read(testBin)

	It("should execute a simple pipeline", func() {
		out1 := new(bytes.Buffer)
		data := NewData()
		cat := &Shell{
			Cmd:  "cat",
			Name: "cat",
		}
		cat.Init()
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p1, err := NewEncoding(
			[]Encoder{cat, cat, cat},
			r,
			data)
		Ω(err).To(BeNil())
		Ω(p1.Exec(out1)).To(BeNil())
		Ω(bytes.Equal(out1.Bytes(), testBin)).To(BeTrue())

		out2 := new(bytes.Buffer)
		zip := &Shell{
			Cmd:  "gzip",
			Name: "zip",
		}
		zip.Init()
		
		r, w = io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p2, err := NewEncoding(
			[]Encoder{zip},
			r,
			data)
		Ω(err).To(BeNil())
		Ω(p2.Exec(out2)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeFalse())

		out3 := new(bytes.Buffer)
		unzip := &Shell{
			Cmd:  "gzip -d",
			Name: "unzip",
		}
		unzip.Init()
		r, w = io.Pipe()
		go func() {
			w.Write(out2.Bytes())
			w.Close()
		}()
		p3, err := NewEncoding(
			[]Encoder{unzip, zip, unzip},
			r,
			data)
		Ω(err).To(BeNil())
		Ω(p3.Exec(out3)).To(BeNil())
		Ω(bytes.Equal(out3.Bytes(), testBin)).To(BeTrue())
	})

	It("should crash with only one command", func() {
		out := new(bytes.Buffer)
		data := NewData()
		crash := &Shell{
			Cmd:  "exit 1",
			Name: "crash",
		}
		crash.Init()
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p, err := NewEncoding(
			[]Encoder{crash},
			r,
			data)
		Ω(err).To(BeNil())
		Ω(p.Exec(out)).ToNot(BeNil())
		Ω(len(out.Bytes())).To(Equal(0))
	})

	It("should crash with pending commands after", func() {
		out := new(bytes.Buffer)
		data := NewData()
		cat := &Shell{
			Cmd:  "cat",
			Name: "cat",
		}
		cat.Init()
		crash := &Shell{
			Cmd:  "exit 1",
			Name: "crash",
		}
		crash.Init()
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p, err := NewEncoding(
			[]Encoder{crash, cat, cat, cat},
			r,
			data)
		Ω(err).To(BeNil())
		Ω(p.Exec(out)).ToNot(BeNil())
		Ω(len(out.Bytes())).To(Equal(0))
	})

	It("should crash with pending encoders before", func() {
		out := new(bytes.Buffer)
		data := NewData()
		cat := &Shell{
			Cmd:  "cat",
			Name: "cat",
		}
		cat.Init()
		crash := &Shell{
			Cmd:  "exit 1",
			Name: "crash",
		}
		crash.Init()
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p, err := NewEncoding(
			[]Encoder{cat, cat, cat, crash, cat, cat},
			r,
			data)
		Ω(err).To(BeNil())
		Ω(p.Exec(out)).ToNot(BeNil())
		Ω(len(out.Bytes())).To(Equal(0))
	})
	
})
