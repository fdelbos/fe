package main

import (
	"github.com/dchest/uniuri"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backend", func() {

	var backend *MongoBackend
	var err error
	It("should create a backend", func() {
		backend, err = NewMongoBackend("bubble", "test-fe", "files")
		Ω(err).To(BeNil())
		Ω(backend).ToNot(BeNil())
	})

	It("should create and id", func() {
		Ω(backend.NewId()).ToNot(BeNil())
		Ω(len(backend.NewId())).ToNot(Equal(0))
	})

	It("should set something and get it", func() {
		id := backend.NewId()
		Ω(id).ToNot(BeNil())

		data := map[string]interface{}{
			"test": uniuri.New(),
		}

		Ω(backend.Set(id, data)).To(BeNil())

		from := make(map[string]interface{})
		from, err := backend.Get(id)
		Ω(err).To(BeNil())
		Ω(from).ToNot(BeNil())
		Ω(from["test"]).ToNot(BeNil())
		Ω(from["test"]).To(Equal(data["test"]))

		Ω(backend.Delete(id)).To(BeNil())
		Ω(backend.Delete(id)).To(Equal(ErrNotFound))

		from = make(map[string]interface{})
		from, err = backend.Get(id)
		Ω(err).To(Equal(ErrNotFound))
	})

})
