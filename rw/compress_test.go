package rw_test

import (
	"bytes"
	"crypto/rand"
	. "github.com/fdelbos/fe/rw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compress", func() {

	testBin := make([]byte, 1<<16)
	rand.Read(testBin)
	out1 := new(bytes.Buffer)
	data := NewData()

	gz := &Gzip{
		Algo: "speed",
	}

	It("should Encode", func() {
		Ω(gz.Init()).To(BeNil())
		Ω(gz.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).To(BeNil())
		Ω(bytes.Equal(out1.Bytes(), testBin)).To(BeFalse())
	})

	out2 := new(bytes.Buffer)
	It("should Decode", func() {
		Ω(gz.Init()).To(BeNil())
		Ω(gz.Decode(
			bytes.NewReader(out1.Bytes()),
			out2,
			data)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())
	})

})
