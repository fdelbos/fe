//
// back.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov 10 2014.
//

package main

import (
	"gopkg.in/mgo.v2"	
	"gopkg.in/mgo.v2/bson"
	// "errors"
	// "mime/multipart"
)
	
var (
	ErrNotFound = mgo.ErrNotFound
)

type Backend interface {
	NewId() string
	Get(string) (map[string]interface{}, error)
	Set(string, map[string]interface{}) error
	Commit(string) error
	Delete(string) error
}

type MongoBackend struct {
	Db *mgo.Database
	Collection string
}

func NewMongoBackend(url, db, collection string) (*MongoBackend, error) {
	b := &MongoBackend{
		Collection: collection,
	}
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	b.Db = session.DB(db)
	return b, nil
}

func (b *MongoBackend) C() *mgo.Collection {
	return b.Db.C(b.Collection)
}

func (b *MongoBackend) NewId() string {
	return bson.NewObjectId().Hex()
}

func (b *MongoBackend) Get(id string) (map[string]interface{}, error) {
	if bson.IsObjectIdHex(id) == false {
		return nil, ErrNotFound
	}
	_id := bson.ObjectIdHex(id)
	m := make(map[string]interface{})
	if err := b.C().FindId(_id).One(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (b *MongoBackend) Set(id string, m map[string]interface{}) error {
	if bson.IsObjectIdHex(id) == false {
		return ErrNotFound
	}
	_id := bson.ObjectIdHex(id)
	m["_id"] = _id
	_, err := b.C().UpsertId(_id, m)
	return err
}

func (b *MongoBackend) Commit(id string) error {
	m, err := b.Get(id)
	if err != nil {
		return err
	}
	m["commit"] = true
	return b.C().UpdateId(m["_id"], m)
}

func (b *MongoBackend) Delete(id string) error {
	if bson.IsObjectIdHex(id) == false {
		return ErrNotFound
	}
	_id := bson.ObjectIdHex(id)	
	return b.C().RemoveId(_id)
}
