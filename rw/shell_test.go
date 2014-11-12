package rw_test

import (
	"bytes"
	"crypto/rand"
	. "github.com/fdelbos/fe/rw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shell", func() {

	It("should Encode", func() {
		testBin := make([]byte, 1<<16)
		rand.Read(testBin)
		out1 := new(bytes.Buffer)
		data := NewData()
		zip := &Shell{
			Cmd:  "gzip",
			Name: "zip",
		}

		Ω(zip.Init()).To(BeNil())
		Ω(zip.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).To(BeNil())
		Ω(bytes.Equal(out1.Bytes(), testBin)).To(BeFalse())

		out2 := new(bytes.Buffer)
		unzip := &Shell{
			Cmd:  "gzip -d",
			Name: "unzip",
		}

		Ω(unzip.Init()).To(BeNil())
		Ω(unzip.Encode(
			bytes.NewReader(out1.Bytes()),
			out2,
			data)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())
	})

	It("should crash", func() {
		testBin := make([]byte, 1<<16)
		rand.Read(testBin)
		out1 := new(bytes.Buffer)
		data := NewData()
		crash := &Shell{
			Cmd:  "exit 1",
			Name: "crash",
		}

		Ω(crash.Init()).To(BeNil())
		Ω(crash.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).ToNot(BeNil())
	})


	
})
