package main

import (
	"github.com/dchest/uniuri"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {

	redis := &RedisCache{
		Host:   "localhost:6379",
		Prefix: "test-fe",
	}
	It("should create a redis cache", func() {
		Ω(redis.Init()).To(BeNil())
	})

	It("should set, get and delete data from redis", func() {
		id := uniuri.New()
		Ω(id).ToNot(BeNil())

		data := struct {
			Str string
		}{"test"}

		Ω(redis.Set(id, data)).To(BeNil())

		from := struct {
			Str string
		}{""}
		Ω(redis.Get(id, &from)).To(BeNil())
		Ω(from.Str).To(Equal("test"))

		Ω(redis.Del(id)).To(BeNil())

		empty := struct {
			Str string
		}{""}
		Ω(redis.Get(id, &empty)).ToNot(BeNil())
	})
})
