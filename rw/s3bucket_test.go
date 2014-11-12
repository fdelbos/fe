package rw_test

import (
	"bytes"
	"crypto/rand"
	"github.com/dchest/uniuri"
	. "github.com/fdelbos/fe/rw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("S3bucket", func() {

	testBin := make([]byte, 1<<16)
	rand.Read(testBin)
	data := NewData()
	s3 := &S3Bucket{
		Bucket: "test-fred",
		Name:   "s3",
	}

	id := uniuri.New()
	data.Set("identifier", id)

	It("should Write", func() {
		Ω(s3.Init()).To(BeNil())
		w, err := s3.NewWriter(data)
		Ω(err).To(BeNil())
		Ω(w).ToNot(BeNil())
		l, err := io.Copy(w, bytes.NewReader(testBin))
		w.Close()
		Ω(err).To(BeNil())
		Ω(len(testBin) == int(l)).To(BeTrue())
	})

	It("should read", func() {
		r, err := s3.NewReader(data)
		Ω(err).To(BeNil())
		Ω(r).ToNot(BeNil())
		out1 := new(bytes.Buffer)
		l, err := io.Copy(out1, r)
		Ω(err).To(BeNil())
		r.Close()
		Ω(len(testBin) == int(l)).To(BeTrue())
		Ω(bytes.Equal(testBin, out1.Bytes())).To(BeTrue())
	})

	It("should delete", func() {
		Ω(s3.Delete(data)).To(BeNil())

		_, err := s3.NewReader(data)
		Ω(err).ToNot(BeNil())
	})

})
