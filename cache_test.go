package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/dchest/uniuri"
)

var _ = Describe("Cache", func() {

	var redis *Redis
	var err error
	It("should create a redis cache", func() {
		redis, err = NewRedis("bubble:6379", "test-fe")
		Ω(err).To(BeNil())
		Ω(redis).ToNot(BeNil())
	})

	It("should set, get and delete data from redis", func(){
		id := uniuri.New()
		Ω(id).ToNot(BeNil())
		
		data := struct {
			Str string
		} {"test"}

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
