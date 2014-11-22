//
// config.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov 13 2014.
//

package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/fdelbos/fe/rw"
	"log"
	"os"
	"regexp"
	"strconv"
	"io/ioutil"
)

type Configuration struct {
	Encoders  map[string]rw.Encoder
	Decoders  map[string]rw.Decoder
	Outputers map[string]rw.Outputer
	Inputers  map[string]rw.Inputer
	Deleters  map[string]rw.Deleter
	Services  map[string]*Service
	Backends  map[string]Backend
	Caches    map[string]Cache
}

func typeToInstance(ty string) (interface{}, error) {
	switch ty {
	case "aes256":
		return new(rw.AES256), nil
	case "gzip":
		return new(rw.Gzip), nil
	case "file":
		return new(rw.File), nil
	case "resize":
		return new(rw.Resize), nil
	case "s3":
		return new(rw.S3Bucket), nil
	case "shell":
		return new(rw.Shell), nil
	default:
		return nil, errors.New("Unsupported module type '" + ty +"'")
	}
	return nil, nil
}

func (c *Configuration) extractRW(name string, raw json.RawMessage) error {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &m); err != nil {
		return err
	}
	if _, exists := m["type"]; exists == false {
		return errors.New("module '" + name + "' doesn't have a type")
	}
	var ty string
	if err := json.Unmarshal(m["type"], &ty); err != nil {
		return err
	}
	obj, err := typeToInstance(ty)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, obj); err != nil {
		return err
	}

	switch v := obj.(type) {
	case rw.Encoder:
		c.Encoders[name] = v
	}
	switch v := obj.(type) {
	case rw.Decoder:
		c.Decoders[name] = v
	}
	switch v := obj.(type) {
	case rw.Outputer:
		c.Outputers[name] = v
	}
	switch v := obj.(type) {
	case rw.Inputer:
		c.Inputers[name] = v
	}
	switch v := obj.(type) {
	case rw.Deleter:
		c.Deleters[name] = v
	}
	return nil
}

func getAuth(right string) Access {
	switch right {
	case "denied":
		return AccDenied
	case "private":
		return AccPrivate
	case "commit":
		return AccCommit
	case "token":
		return AccToken
	case "public":
		return AccPublic
	default:
		log.Fatal("Unknow access right: " + right)
	}
	return AccDenied
}

type JsonService struct {
	Cache   string `json:"cache"`
	Backend string `json:"backend"`
	Get     struct {
		Access   string   `json:"access"`
		Pipeline []string `json:"pipeline"`
		Input    string   `json:"input"`
	} `json:"get,omitempty"`
	Post struct {
		Access   string   `json:"access"`
		Pipeline []string `json:"pipeline,omitempty"`
		Output   string   `json:"output"`
	} `json:"post,omitempty"`
	Delete struct {
		Access  string `json:"access"`
		Deleter string `json:"deleter"`
	} `json:"delete,omitempty"`
}

func (c *Configuration) buildServicePost(s *Service, js *JsonService) error {
	s.Post = getAuth(js.Post.Access)
	if js.Post.Output == "" {
		return errors.New("POST method requires an 'output'")
	}
	if _, exists := c.Outputers[js.Post.Output]; exists == false {
		return errors.New("output '" + js.Post.Output + "' is undefined")
	}
	s.EncodingPipe = &rw.EncodingPipeline{
		Encoders: make([]rw.Encoder, 0),
		Output:   c.Outputers[js.Post.Output],
	}
	if js.Post.Pipeline != nil && len(js.Post.Pipeline) > 0 {
		for _, v := range js.Post.Pipeline {
			if _, exists := c.Encoders[v]; exists == false {
				log.Fatal("encoder '" + v + "' is undefined")
			}
			s.EncodingPipe.Encoders = append(s.EncodingPipe.Encoders, c.Encoders[v])
		}
	}
	return nil
}

func (c *Configuration) buildServiceGet(s *Service, js *JsonService) error {
	s.Get = getAuth(js.Get.Access)
	if s.Get != AccPrivate && s.Get != AccPublic {
		return errors.New("GET access can only be 'denied', 'private' or 'public'")
	}
	s.DecodingPipe = &rw.DecodingPipeline{
		Decoders: make([]rw.Decoder, 0),
		Input:    nil,
	}
	if js.Get.Input == "" {
		if s.EncodingPipe != nil {
			switch v := s.EncodingPipe.Output.(type) {
			case rw.Inputer:
				s.DecodingPipe.Input = v
			default:
				return errors.New("GET method requires and 'input'")
			}
		} else {
			return errors.New("GET method requires and 'input'")
		}
	} else {
		if _, exists := c.Inputers[js.Get.Input]; exists == false {
			return errors.New("input '" + js.Get.Input + "' is undefined")
		}
		s.DecodingPipe.Input = c.Inputers[js.Get.Input]
	}

	if js.Get.Pipeline != nil && len(js.Get.Pipeline) > 0 {
		for _, v := range js.Get.Pipeline {
			if _, exists := c.Decoders[v]; exists == false {
				return errors.New("decoder '" + v + "' is undefined")
			}
			s.DecodingPipe.Decoders = append(s.DecodingPipe.Decoders, c.Decoders[v])
		}
	} else {
		for i := len(s.EncodingPipe.Encoders) - 1; i >= 0; i -= 1 {
			switch v := s.EncodingPipe.Encoders[i].(type) {
			case rw.Decoder:
				s.DecodingPipe.Decoders = append(s.DecodingPipe.Decoders, v)
			}
		}
	}
	return nil
}

var validUrl = regexp.MustCompile(`^(/[^/]+)+$`)
func (c *Configuration) extractService(url string, raw json.RawMessage) error {
	if validUrl.MatchString(url) == false {
		return errors.New(" '" + url + "' is not a invalid url")
	}

	js := new(JsonService)
	if err := json.Unmarshal(raw, js); err != nil {
		return err
	}
	s := &Service{
		Url:    url,
		Post:   AccDenied,
		Get:    AccDenied,
		Delete: AccDenied,
	}
	if js.Backend != "" {
		if _, exists := c.Backends[js.Backend]; exists == false {
			return errors.New("service '" + url + "': bakend  '" + js.Backend + "' not found")
		}
		s.Back = c.Backends[js.Backend]
	}
	if js.Cache != "" {
		if _, exists := c.Caches[js.Cache]; exists == false {
			return errors.New("cache '" + js.Cache + "' doesn't exists")
		}
		tokens := &TokenService{
			Service: url,
			cache:   c.Caches[js.Cache],
		}
		s.Tokens = tokens
	}
	if js.Post.Access != "" && js.Post.Access != "denied" {
		if err := c.buildServicePost(s, js); err != nil {
			return err
		}
	}
	if js.Get.Access != "" && js.Get.Access != "denied" {
		if err := c.buildServiceGet(s, js); err != nil {
			return err
		}
	}
	// need to do delete
	c.Services[url] = s
	return nil
}

func (c *Configuration) extractBackends(definition map[string]json.RawMessage) error {
	for k,v := range definition {
		m := make(map[string]json.RawMessage)
		if err := json.Unmarshal(v, &m); err != nil {
			return err
		}
		var ty string
		if err := json.Unmarshal(m["type"], &ty); err != nil {
			return err
		}
		if ty != "mongodb" {
			return errors.New("Unsupported backend '" + ty + "'")
		}
		b := new(MongoBackend)
		if err := json.Unmarshal(v, b); err != nil {
			return err
		}
		if err := b.Init(); err != nil {
			return err
		}
		c.Backends[k] = b
	}
	return nil
}

func (c *Configuration) extractCaches(definition map[string]json.RawMessage) error {
	for k,v := range definition {
		m := make(map[string]json.RawMessage)
		if err := json.Unmarshal(v, &m); err != nil {
			return err
		}
		var ty string
		if err := json.Unmarshal(m["type"], &ty); err != nil {
			return err
		}
		if ty != "redis" {
			return errors.New("Unsupported backend '" + ty + "'")
		}
		b := new(RedisCache)
		if err := json.Unmarshal(v, b); err != nil {
			return err
		}
		if err := b.Init(); err != nil {
			return err
		}
		c.Caches[k] = b
	}
	return nil
}

func newConfiguration(js []byte) (*Configuration, error) {
	c := &Configuration{
		Encoders:  make(map[string]rw.Encoder),
		Decoders:  make(map[string]rw.Decoder),
		Outputers: make(map[string]rw.Outputer),
		Inputers:  make(map[string]rw.Inputer),
		Deleters:  make(map[string]rw.Deleter),
		Services:  make(map[string]*Service),
		Backends:  make(map[string]Backend),
		Caches:    make(map[string]Cache),
	}

	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(js, &m); err != nil {
		return nil, err
	}
	if _, exists := m["caches"]; exists == true {
		caches := make(map[string]json.RawMessage)
		if err := json.Unmarshal(m["caches"], &caches); err != nil {
			return nil, err
		}
		if err := c.extractCaches(caches); err != nil {
			return nil, err
		}
	}
	if _, exists := m["backends"]; exists == true {
		backends := make(map[string]json.RawMessage)
		if err := json.Unmarshal(m["backends"], &backends); err != nil {
			return nil, err
		}
		if err := c.extractBackends(backends); err != nil {
			return nil, err
		}
	}
	
	if _, exists := m["modules"]; exists == true {
		modules := make(map[string]json.RawMessage)
		if err := json.Unmarshal(m["modules"], &modules); err != nil {
			return nil, err
		}
		for k, v := range modules {
			if err := c.extractRW(k, v); err != nil {
				return nil, err
			}
		}
	}
	if _, exists := m["api"]; exists == true {
		api := make(map[string]json.RawMessage)
		if err := json.Unmarshal(m["api"], &api); err != nil {
			return nil, err
		}
		for k, v := range api {
			if err := c.extractService(k, v); err != nil {
				return nil, err
			}
		}
	}
	fmt.Println(c)
	return c, nil
}

func genAes(size int) {
	var buff []byte
	switch size {
	case 128:
		buff = make([]byte, 16)
	case 256:
		buff = make([]byte, 32)
	default:
		log.Fatal("Key size must be 256 for AES 256 or 128 for AES 128")
	}
	if _, err := rand.Read(buff); err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(buff)
	os.Exit(0)
}

func parseConfig() (*Configuration, error) {

	msg := `fe - File Exchange proxy

Usage:
  fe [--daemon --private <port> --public <port>] <config>
  fe genaes <keysize>
  fe -h | --help
  fe -v | --version

Options:
  -h, --help              Show this screen
  -v, --version           Show version
  -d, --daemon            Run as a daemon
  -p, --private <port>    Set Private port  [default: 7031]
  -u, --public <port>     Set Public port   [default: 7032]
`
	config, err := docopt.Parse(msg, nil, true, "fe 0.1.1", false)
	if err != nil {
		return nil, err
	}
	fmt.Println(config)
	if config["genaes"].(bool) == true {
		size := config["<keysize>"].(string)
		i, err := strconv.Atoi(size)
		if err != nil {
			return nil, err
		}
		genAes(i)
	}
	name, exists := config["<config>"]
	if exists == false {
		return nil, errors.New("No configuration file provided")
	}
	js, err := ioutil.ReadFile(name.(string))
	if err != nil {
		return nil, err
	}
	return newConfiguration(js)
}
