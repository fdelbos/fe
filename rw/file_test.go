package rw_test

import (
	. "github.com/fdelbos/fe/rw"

	"bytes"
	"crypto/rand"
	"github.com/dchest/uniuri"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"os"
)

var _ = Describe("File", func() {

	testBin := make([]byte, 1<<16)
	rand.Read(testBin)
	data := NewData()
	f := &File{
		Dir:  "/tmp/" + uniuri.New(),
		Name: "file",
	}

	id := uniuri.New()
	data.Set("identifier", id)

	It("should Write", func() {
		Ω(f.Init()).To(BeNil())
		w, err := f.NewWriter(data)
		Ω(err).To(BeNil())
		Ω(w).ToNot(BeNil())
		l, err := io.Copy(w, bytes.NewReader(testBin))
		w.Close()
		Ω(err).To(BeNil())
		Ω(len(testBin) == int(l)).To(BeTrue())
	})

	It("should read", func() {
		r, err := f.NewReader(data)
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
		Ω(f.Delete(data)).To(BeNil())
		_, err := os.Stat(f.Dir + "/" + id)
		Ω(os.IsNotExist(err)).To(BeTrue())
	})
})
