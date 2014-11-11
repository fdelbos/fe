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
			return RwError(c, err.Error())
		}
	case c.Base64String != "":
		key, err = base64.StdEncoding.DecodeString(c.Base64String)
		if err != nil {
			return RwError(c, err.Error())
		}
	default:
		return RwError(c, "needs a key")
	}
	if len(key) != 32 {
		return RwError(c, "key must be 32 bytes long")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return RwError(c, err.Error())
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
		return RwError(c, err.Error())
	}
	return nil
}

func (c *AES256) Decode(r io.Reader, w io.Writer, d *Data) error {
	iv64 := d.Get("iv")
	if iv64 == "" {
		return RwError(c, "AES 256: no initialization vector provided")
	}
	iv, err := base64.StdEncoding.DecodeString(iv64.(string))
	if err != nil {
		return RwError(c, err.Error())
	}
	stream := cipher.NewCFBDecrypter(c.block, iv)
	reader := &cipher.StreamReader{S: stream, R: r}
	if _, err := io.Copy(w, reader); err != nil {
		return RwError(c, err.Error())
	}
	return nil
}

func generateIV(bytes int) []byte {
	b := make([]byte, bytes)
	rand.Read(b)
	return b
}
