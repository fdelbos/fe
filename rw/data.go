//
// data.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  9 2014.
//

package rw

import (
	"errors"
	"sync"
)

type Data struct {
	sync.RWMutex
	data map[string]interface{}
}

func NewData() *Data {
	return &Data{
		data: make(map[string]interface{}),
	}
}

func (d *Data) Get(key string) (interface{}, error) {
	d.RLock()
	defer d.RUnlock()
	v, exists := d.data[key]
	if exists == false {
		return nil, errors.New("Data '" + key + "' not found")
	}
	return v, nil
}

func (d *Data) Set(key string, value interface{}) {
	d.Lock()
	defer d.Unlock()
	d.data[key] = value
}
