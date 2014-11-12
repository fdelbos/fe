package rw_test

import (
	"bytes"
	"crypto/rand"
	. "github.com/fdelbos/fe/rw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("Pipeline", func() {

	testBin := make([]byte, 1<<18)
	rand.Read(testBin)

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

	It("should execute a simple pipeline", func() {
		out1 := new(bytes.Buffer)
		data := NewData()
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p1 := NewEncoding(
			[]Encoder{cat, cat, cat},
			r,
			data)
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
		p2 := NewEncoding(
			[]Encoder{zip},
			r,
			data)
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
		p3 := NewEncoding(
			[]Encoder{unzip, zip, unzip},
			r,
			data)
		Ω(p3.Exec(out3)).To(BeNil())
		Ω(bytes.Equal(out3.Bytes(), testBin)).To(BeTrue())
	})

	It("should crash with only one command", func() {
		out := new(bytes.Buffer)
		data := NewData()
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p := NewEncoding(
			[]Encoder{crash},
			r,
			data)
		Ω(p.Exec(out)).ToNot(BeNil())
		Ω(len(out.Bytes())).To(Equal(0))
	})

	It("should crash with pending commands after", func() {
		out := new(bytes.Buffer)
		data := NewData()
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p := NewEncoding(
			[]Encoder{crash, cat, cat, cat},
			r,
			data)
		Ω(p.Exec(out)).ToNot(BeNil())
		Ω(len(out.Bytes())).To(Equal(0))
	})

	It("should crash with pending encoders before", func() {
		out := new(bytes.Buffer)
		data := NewData()
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p := NewEncoding(
			[]Encoder{cat, cat, cat, crash, cat, cat},
			r,
			data)
		Ω(p.Exec(out)).ToNot(BeNil())
		Ω(len(out.Bytes())).To(Equal(0))
	})

	It("should do multiple crash", func() {
		out := new(bytes.Buffer)
		data := NewData()
		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p := NewEncoding(
			[]Encoder{cat, cat, crash, cat, crash, cat, crash, cat},
			r,
			data)
		Ω(p.Exec(out)).ToNot(BeNil())
		Ω(len(out.Bytes())).To(Equal(0))
	})

	It("should do encode / decode", func() {
		out := new(bytes.Buffer)
		data := NewData()
		aes := &AES256{
			Base64String: "ETl5QyPnHfi+vF4HrZfFvO2Julv4LVL7HNB1N7vkLGU=",
			Name:         "aes",
		}
		aes.Init()

		r, w := io.Pipe()
		go func() {
			w.Write(testBin)
			w.Close()
		}()
		p1 := NewEncoding(
			[]Encoder{aes},
			r,
			data)
		Ω(p1.Exec(out)).To(BeNil())
		Ω(len(out.Bytes())).ToNot(Equal(0))
		Ω(bytes.Equal(out.Bytes(), testBin)).To(BeFalse())

		r2, w2 := io.Pipe()
		go func() {
			w2.Write(out.Bytes())
			w2.Close()
		}()
		out2 := new(bytes.Buffer)

		p2 := NewDecoding(
			[]Decoder{aes},
			r2,
			data)
		Ω(p2.Exec(out2)).To(BeNil())
		Ω(len(out2.Bytes())).ToNot(Equal(0))
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())
	})

})
