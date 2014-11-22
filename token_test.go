package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Token", func() {

	var ts *TokenService
	redis := &RedisCache{
		Host:   "localhost:6379",
		Prefix: "test-fe",
	}
	It("should create a redis cache and a token service", func() {
		Ω(redis.Init()).To(BeNil())
		ts = &TokenService{
			Service: "test",
			cache:   redis,
		}
	})

	It("should create, get, set, delete tokens", func() {
		t1, err := ts.NewToken("test")
		Ω(err).To(BeNil())
		Ω(t1).ToNot(BeNil())

		t2, err := ts.Get(t1.Key)
		Ω(err).To(BeNil())
		Ω(t2).ToNot(BeNil())
		Ω(t2.Key).To(Equal(t1.Key))

		t2.Identifier = "test"
		Ω(ts.Set(t2)).To(BeNil())

		t3, err := ts.Get(t1.Key)
		Ω(err).To(BeNil())
		Ω(t3).ToNot(BeNil())
		Ω(t3.Key).To(Equal(t1.Key))
		Ω(t3.Identifier).To(Equal("test"))

		Ω(ts.Del(t3.Key)).To(BeNil())
		_, err = ts.Get(t1.Key)
		Ω(err).ToNot(BeNil())
	})

})
