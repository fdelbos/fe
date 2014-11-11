package rw_test

import (
	"bytes"
	"crypto/rand"
	. "github.com/fdelbos/fe/rw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aes", func() {

	testBin := make([]byte, 1<<16)
	rand.Read(testBin)
	out1 := new(bytes.Buffer)
	data := NewData()

	aes := &AES256{
		Base64String: "ETl5QyPnHfi+vF4HrZfFvO2Julv4LVL7HNB1N7vkLGU=",
		Name:         "aes",
	}

	It("should Encode", func() {
		Ω(aes.Init()).To(BeNil())
		Ω(aes.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).To(BeNil())
		Ω(bytes.Equal(out1.Bytes(), testBin)).To(BeFalse())
		Ω(len(data.Get("iv").(string)) > 0).To(BeTrue())
	})

	out2 := new(bytes.Buffer)
	It("should Decode", func() {
		Ω(aes.Init()).To(BeNil())
		Ω(aes.Decode(
			bytes.NewReader(out1.Bytes()),
			out2,
			data)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())
	})
})
