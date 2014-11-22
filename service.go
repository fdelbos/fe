//
// service.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  9 2014.
//

package main

import (
	"errors"
	"github.com/fdelbos/fe/rw"
	"github.com/gorilla/mux"
	"io"
)

type Access int

var NoEncodindError = errors.New("No encoding pipeline defined")
var NoDecodindError = errors.New("No decoding pipeline defined")

const (
	AccDenied Access = iota
	AccPrivate
	AccCommit
	AccToken
	AccPublic
)

type Service struct {
	Url          string
	Back         Backend
	Post         Access
	Get          Access
	Delete       Access
	EncodingPipe *rw.EncodingPipeline
	DecodingPipe *rw.DecodingPipeline
	Tokens       *TokenService
}

func (s *Service) Encode(id string, r io.ReadCloser, data *rw.Data) error {
	defer data.Set("identifier", id)

	if s.EncodingPipe == nil {
		return NoEncodindError
	}
	w, err := s.EncodingPipe.Output.NewWriter(id, data)
	if err != nil {
		return err
	}
	if len(s.EncodingPipe.Encoders) == 0 {
		l, err := io.Copy(w, r)
		if err != nil {
			return err
		}
		w.Close()
		data.Set("size", l)
	} else {
		p := rw.NewEncoding(s.EncodingPipe.Encoders, r, data)
		if err := p.Exec(w); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Decode(id string, w io.WriteCloser, data *rw.Data) error {
	if s.DecodingPipe == nil {
		return NoDecodindError
	}
	r, err := s.DecodingPipe.Input.NewReader(id, data)
	if err != nil {
		return err
	}
	if len(s.DecodingPipe.Decoders) == 0 {
		done := make(chan error)

		go func() {
			defer r.Close()
			_, err := io.Copy(w, r)
			done <- err
		}()
		err := <-done
		if err != nil {
			return err
		}
	} else {
		p := rw.NewDecoding(s.DecodingPipe.Decoders, r, data)
		if err := p.Exec(w); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) RegisterPublic(r *mux.Router) {
	sr := r.PathPrefix(s.Url).Subrouter()

	if s.Post == AccPublic {
		sr.HandleFunc("/", s.publicPost).Methods("POST")
	} else if s.Post == AccToken || s.Post == AccCommit {
		sr.HandleFunc("/{token}", s.publicPost).Methods("POST")
	}
	if s.Get == AccPublic {
		sr.HandleFunc("/{identifier}", s.publicGet).Methods("GET")
	}
}

func (s *Service) setCommit(r *mux.Router) {
	sr := r.PathPrefix("/commit").Subrouter()
	sr.HandleFunc("/", s.commit).Methods("POST")
}

func (s *Service) setTokens(r *mux.Router) {
	sr := r.PathPrefix("/token").Subrouter()
	sr.HandleFunc("/", s.genToken).Methods("POST")
	sr.HandleFunc("/", s.genToken).Methods("GET")
	sr.HandleFunc("/{token}", s.getToken).Methods("GET")
	sr.HandleFunc("/{token}", s.deleteToken).Methods("DELETE")
}

func (s *Service) RegisterPrivate(r *mux.Router) {
	sr := r.PathPrefix(s.Url).Subrouter()

	if s.Post != AccDenied {
		if s.Post == AccCommit {
			s.setCommit(sr)
		}
		if s.Post == AccCommit || s.Post == AccToken {
			s.setTokens(sr)
		}
		sr.HandleFunc("/", s.privatePost).Methods("POST")
	}
	if s.Get != AccDenied {
		sr.HandleFunc("/{identifier}", s.privateGet).Methods("GET")
	}
	if s.Delete != AccDenied {
		sr.HandleFunc("/{identifier}", s.privateDelete).Methods("DELETE")
	}

}
