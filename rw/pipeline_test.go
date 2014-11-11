package rw_test

import (
	"bytes"
	"crypto/rand"
	. "github.com/fdelbos/fe/rw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipeline", func() {

	testBin := make([]byte, 1<<16)
	rand.Read(testBin)

	It("should execute a simple pipeline", func() {
		out1 := new(bytes.Buffer)
		data := NewData()
		cat := &Shell{
			Cmd:  "cat",
			Name: "cat",
		}
		cat.Init()
		p1, err := NewEncoding(
			[]Encoder{cat, cat, cat},
			bytes.NewReader(testBin),
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
		p2, err := NewEncoding(
			[]Encoder{zip},
			bytes.NewReader(testBin),
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
		p3, err := NewEncoding(
			[]Encoder{unzip, zip, unzip},
			bytes.NewReader(out2.Bytes()),
			data)
		Ω(err).To(BeNil())
		Ω(p3.Exec(out3)).To(BeNil())
		Ω(bytes.Equal(out3.Bytes(), testBin)).To(BeTrue())
	})
})
