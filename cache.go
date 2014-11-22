//
// redis.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov 10 2014.
// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
//

package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var ErrCacheNotFound = redis.ErrNil

type Cache interface {
	Set(string, interface{}) error
	Get(string, interface{}) error
	Del(string) error
	Init() error
}

type RedisCache struct {
	Prefix string `json:"prefix"`
	Host   string `json:"host"`
	pool   *redis.Pool
}

func (r *RedisCache) Init() error {
	r.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", r.Host)
			if err != nil {
				log.Fatal(err)
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func (r *RedisCache) Set(key string, data interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", r.Prefix+key, j)
	return err
}

func (r *RedisCache) Get(key string, container interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()
	data, err := conn.Do("GET", r.Prefix+key)
	if err != nil {
		return err
	}
	if data == nil {
		return ErrCacheNotFound
	}
	return json.Unmarshal(data.([]byte), container)
}

func (r *RedisCache) Del(key string) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", r.Prefix+key)
	return err
}
