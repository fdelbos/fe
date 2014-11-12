//
// aes.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  8 2014.
//

package rw

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
)

type AES256 struct {
	block        cipher.Block
	File         string `json:"file,omitempty"`
	Base64String string `json:"base64,omitempty"`
	Name         string `json:"-"`
}

func (c *AES256) GetName() string {
	return c.Name
}

func (c *AES256) Init() error {
	var key []byte
	var err error

	switch {
	case c.File != "":
		key, err = ioutil.ReadFile(c.File)
		if err != nil {
			return err
		}
	case c.Base64String != "":
		key, err = base64.StdEncoding.DecodeString(c.Base64String)
		if err != nil {
			return err
		}
	default:
		return errors.New("needs a cycpher key")
	}
	if len(key) != 32 {
		return errors.New("key must be 32 bytes long")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	c.block = block
	return nil
}

func (c *AES256) Encode(r io.Reader, w io.Writer, d *Data) error {
	iv := generateIV(c.block.BlockSize())
	d.Set("iv", base64.StdEncoding.EncodeToString(iv))

	stream := cipher.NewCFBEncrypter(c.block, iv)
	writer := &cipher.StreamWriter{S: stream, W: w}
	if _, err := io.Copy(writer, r); err != nil {
		return err
	}
	return nil
}

func (c *AES256) Decode(r io.Reader, w io.Writer, d *Data) error {
	iv64, err := d.Get("iv")
	if err != nil {
		return err
	}
	if iv64.(string) == "" {
		return errors.New("no initialization vector provided")
	}
	iv, err := base64.StdEncoding.DecodeString(iv64.(string))
	if err != nil {
		return err
	}
	stream := cipher.NewCFBDecrypter(c.block, iv)
	reader := &cipher.StreamReader{S: stream, R: r}
	if _, err := io.Copy(w, reader); err != nil {
		return err
	}
	return nil
}

func generateIV(bytes int) []byte {
	b := make([]byte, bytes)
	rand.Read(b)
	return b
}
