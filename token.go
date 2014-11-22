//
// token.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov 10 2014.
//

package main

import (
	"github.com/dchest/uniuri"
)

type Token struct {
	Key        string `json:"key"`
	Service    string `json:"service"`
	Identifier string `json:"identifier,omitempty"`
}

type TokenService struct {
	Service string
	cache   Cache
}

func (ts *TokenService) NewToken(service string) (*Token, error) {
	token := &Token{
		Key:        uniuri.New(),
		Service:    service,
		Identifier: "",
	}
	if err := ts.Set(token); err != nil {
		return nil, err
	}
	return token, nil
}

func (ts *TokenService) Get(key string) (*Token, error) {
	token := new(Token)
	if err := ts.cache.Get(key, token); err != nil {
		return nil, err
	}
	return token, nil
}

func (ts *TokenService) Del(key string) error {
	return ts.cache.Del(key)
}

func (ts *TokenService) Set(token *Token) error {
	return ts.cache.Set(token.Key, token)
}
