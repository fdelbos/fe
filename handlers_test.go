package main

import (
	"bytes"
	"crypto/rand"
	"github.com/dchest/uniuri"
	"github.com/fdelbos/fe/rw"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	//	"io"
	//	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/url"
)

var _ = Describe("Handlers", func() {

	testBin := make([]byte, 1<<16)
	rand.Read(testBin)

	file := &rw.File{
		Dir:  "/tmp/" + uniuri.New(),
		Name: "file",
	}
	file.Init()

	aes := &rw.AES256{
		Base64String: "ETl5QyPnHfi+vF4HrZfFvO2Julv4LVL7HNB1N7vkLGU=",
		Name:         "aes",
	}
	aes.Init()

	service := &Service{
		Url: "/test",
		EncodingPipe: &rw.EncodingPipeline{
			Encoders: []rw.Encoder{aes},
			Output:   file,
		},
		DecodingPipe: &rw.DecodingPipeline{
			Decoders: []rw.Decoder{aes},
			Input:    file,
		},
		Post:   AccCommit,
		Get:    AccPublic,
		Delete: AccPrivate,
	}

	var backend *MongoBackend
	var redis *Redis
	var ts *httptest.Server

	It("should start backend, cache", func() {
		var err error
		backend, err = NewMongoBackend("bubble", "test-fe", "files")
		Ω(err).To(BeNil())
		Ω(backend).ToNot(BeNil())
		service.Back = backend

		redis, err = NewRedis("bubble:6379", "test-fe")
		Ω(err).To(BeNil())
		Ω(redis).ToNot(BeNil())
		service.Tokens = &TokenService{
			Service: "test",
			cache:   redis,
		}
		r := mux.NewRouter()
		service.RegisterPrivate(r.PathPrefix("/private").Subrouter())
		service.RegisterPublic(r.PathPrefix("/public").Subrouter())
		ts = httptest.NewServer(r)

	})

	var id string

	It("should post to privatePost", func() {

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("my-file", "file.test")
		part.Write(testBin)
		Ω(err).To(BeNil())
		writer.Close()

		req, err := http.NewRequest("POST", ts.URL+"/private/test/", body)
		Ω(err).To(BeNil())
		req.Header.Add("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(201))
		Ω(resp.Header.Get("Content-Type")).To(Equal("application/json"))
		buff := new(bytes.Buffer)
		buff.ReadFrom(resp.Body)
		m := make(map[string]interface{})

		Ω(json.Unmarshal(buff.Bytes(), &m)).To(BeNil())
		Ω(m["identifier"].(string)).ToNot(BeNil())
		Ω(m["commit"]).To(BeNil())
		Ω(m["iv"]).To(BeNil())
		id = m["identifier"].(string)
	})

	It("should get privateGet", func() {
		resp, err := http.Get(ts.URL + "/private/test/" + id)
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(200))
		buff := new(bytes.Buffer)
		buff.ReadFrom(resp.Body)
		Ω(bytes.Equal(buff.Bytes(), testBin)).To(BeTrue())

		resp, err = http.Get(ts.URL + "/private" + service.Url + "/wrong")
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(404))
	})

	var tokenKey string

	It("should get genToken", func() {
		resp, err := http.Get(ts.URL + "/private/test/token/")
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(201))

		buff := new(bytes.Buffer)
		buff.ReadFrom(resp.Body)

		m := make(map[string]interface{})
		Ω(json.Unmarshal(buff.Bytes(), &m)).To(BeNil())
		Ω(m["service"].(string)).To(Equal(service.Url))
		Ω(m["key"].(string)).ToNot(Equal(""))
		tokenKey = m["key"].(string)
	})

	It("should getToken", func() {
		resp, err := http.Get(ts.URL + "/private/test/token/" + tokenKey)
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(200))

		buff := new(bytes.Buffer)
		buff.ReadFrom(resp.Body)

		m := make(map[string]interface{})
		Ω(json.Unmarshal(buff.Bytes(), &m)).To(BeNil())
		Ω(m["service"].(string)).To(Equal(service.Url))
		Ω(m["key"].(string)).To(Equal(tokenKey))
	})

	It("should deleteToken", func() {
		req, err := http.NewRequest("DELETE", ts.URL+"/private/test/token/"+tokenKey, nil)
		Ω(err).To(BeNil())
		resp, err := http.DefaultClient.Do(req)
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(204))

		resp2, err := http.Get(ts.URL + "/private/test/token/" + tokenKey)
		Ω(err).To(BeNil())
		Ω(resp2.StatusCode).To(Equal(404))
	})

	It("get a token again", func() {
		resp, err := http.Get(ts.URL + "/private/test/token/")
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(201))

		buff := new(bytes.Buffer)
		buff.ReadFrom(resp.Body)

		m := make(map[string]interface{})
		Ω(json.Unmarshal(buff.Bytes(), &m)).To(BeNil())
		Ω(m["service"].(string)).To(Equal(service.Url))
		Ω(m["key"].(string)).ToNot(Equal(""))
		tokenKey = m["key"].(string)
	})

	It("should post to publicPost", func() {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("my-public-file", "file.public.test")
		part.Write(testBin)
		Ω(err).To(BeNil())
		writer.Close()

		req, err := http.NewRequest("POST", ts.URL+"/public/test/"+tokenKey, body)
		Ω(err).To(BeNil())
		req.Header.Add("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(201))
		Ω(resp.Header.Get("Content-Type")).To(Equal("application/json"))
		buff := new(bytes.Buffer)
		buff.ReadFrom(resp.Body)
		m := make(map[string]interface{})

		Ω(json.Unmarshal(buff.Bytes(), &m)).To(BeNil())
		Ω(m["identifier"].(string)).ToNot(BeNil())
		Ω(m["commit"]).To(BeNil())
		Ω(m["iv"]).To(BeNil())
		id = m["identifier"].(string)
	})

	It("should get publicGet", func() {
		resp, err := http.Get(ts.URL + "/public/test/" + id)
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(200))
		buff := new(bytes.Buffer)
		buff.ReadFrom(resp.Body)
		Ω(bytes.Equal(buff.Bytes(), testBin)).To(BeTrue())

		resp, err = http.Get(ts.URL + "/public/test/" + tokenKey)
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(404))
	})

	It("should commit", func() {
		resp, err := http.PostForm(ts.URL+"/private/test/commit/",
			url.Values{"token": {tokenKey}})
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(201))
	})
})
