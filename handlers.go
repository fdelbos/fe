//
// handlers.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov 10 2014.
//

package main

import (
	"encoding/json"
	"github.com/fdelbos/fe/rw"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

func (s *Service) decodeMultipartPost(w http.ResponseWriter, r *http.Request, t *Token) {
	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		filename := part.FileName()
		if filename == "" {
			continue
		} else {
			id := s.Back.NewId()
			data := rw.NewData()
			data.Set("identifier", id)
			data.Set("filename", filename)

			if err = s.Encode(id, part, data); err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}

			if t != nil {
				t.Identifier = id
				if err := s.Tokens.Set(t); err != nil {
					http.Error(w, http.StatusText(500), 500)
					return
				}
				if s.Post == AccCommit {
					data.Set("commit", false)
				}
			}

			if err := s.Back.Set(id, data.Export()); err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}

			js, err := data.Filter()
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(201)
			w.Write(js)
			return
		}
	}
	http.Error(w, http.StatusText(400), 400)
}

func (s *Service) privatePost(w http.ResponseWriter, r *http.Request) {
	s.decodeMultipartPost(w, r, nil)
}

func (s *Service) privateGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["identifier"]
	if exists == false {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	m, err := s.Back.Get(id)
	if err != nil {
		if err == ErrNotFound {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}
	data := rw.NewDataFrom(m)
	pr, pw := io.Pipe()

	go io.Copy(w, pr)

	if err = s.Decode(id, pw, data); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func (s *Service) privateDelete(w http.ResponseWriter, r *http.Request) {
}

func (s *Service) publicPost(w http.ResponseWriter, r *http.Request) {
	var token *Token
	token = nil
	if s.Post == AccToken || s.Post == AccCommit {
		vars := mux.Vars(r)
		id, exists := vars["token"]
		if exists == false {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		var err error
		token, err = s.Tokens.Get(id)
		if err == ErrCacheNotFound {
			http.Error(w, http.StatusText(403), 403)
			return
		}
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if token.Service != s.Url || token.Identifier != "" {
			http.Error(w, http.StatusText(403), 403)
			return
		}
	}
	if token == nil && s.Post != AccPublic {
		http.Error(w, http.StatusText(403), 403)
		return
	}
	s.decodeMultipartPost(w, r, token)
}

func (s *Service) publicGet(w http.ResponseWriter, r *http.Request) {
	s.privateGet(w, r)
}

func (s *Service) commit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	token := r.FormValue("token")
	if token == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	t, err := s.Tokens.Get(token)
	if err == ErrCacheNotFound {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if t.Identifier == "" {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	m, err := s.Back.Get(t.Identifier)
	if err != nil {
		if err == ErrNotFound {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if _, exists := m["commit"]; exists == false {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	c := m["commit"]
	switch c := c.(type) {
	case bool:
		if c != false {
			http.Error(w, http.StatusText(404), 404)
			return
		}
	default:
		http.Error(w, http.StatusText(404), 404)
		return
	}

	m["commitDate"] = time.Now()
	m["commit"] = true
	if err := s.Back.Set(t.Identifier, m); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	d := map[string]interface{}{
		"size":       m["size"],
		"identifier": m["identifier"],
		"filename":   m["filename"],
		"commit":     m["commit"],
		"commitDate": m["commitDate"],
	}
	js, err := json.Marshal(d)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(js)
	return
}

func (s *Service) genToken(w http.ResponseWriter, r *http.Request) {
	t, err := s.Tokens.NewToken(s.Url)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	js, err := json.Marshal(t)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(js)
	return
}

func (s *Service) getToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["token"]
	if exists == false {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	t, err := s.Tokens.Get(id)
	if err == ErrCacheNotFound {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	js, err := json.Marshal(t)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
	return
}

func (s *Service) deleteToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["token"]
	if exists == false {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	err := s.Tokens.Del(id)
	if err == ErrCacheNotFound {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	http.Error(w, http.StatusText(204), 204)
}
