//
// token.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov 10 2014.
// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
//

package main

import (
	"github.com/dchest/uniuri"
)

type TokenOperation int

const (
	TokPost   TokenOperation = iota
	TokGet    TokenOperation = iota
	TokUpdate TokenOperation = iota
	TokDelete TokenOperation = iota
)

type Token struct {
	Key        string         `json:"key"`
	Service    string         `json:"service"`
	Operation  TokenOperation `json:"operation"`
	Identifier string         `json:"identifier,omitempty"`
}

type TokenService struct {
	Service string
	cache   Cache
}

func (ts *TokenService) NewToken(service string, operation TokenOperation, identifier string) (*Token, error) {
	token := &Token{
		Key:        uniuri.New(),
		Service:    service,
		Operation:  operation,
		Identifier: identifier,
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
