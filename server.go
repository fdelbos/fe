//
// server.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  8 2014.
//

package main

import (
	"fmt"
	"log"
	"github.com/fdelbos/fe/rw"
	"net/http"
	"crypto/rand"
	"encoding/base64"
)

func GenCypher(cypher string) {
	var buff []byte
	switch cypher {
	case "aes128":
		buff = make([]byte, 16)
	case "aes256":
		buff = make([]byte, 32)
	default:
		log.Fatalf("cypher '%s' not supported!", cypher)
	}
	if _, err := rand.Read(buff); err != nil {
		log.Fatal(err)
	}
	fmt.Print(base64.StdEncoding.EncodeToString(buff))
}

func uploadHandler2(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Println("1")
		fmt.Fprintln(w, err)
		return
	}

	fmt.Println(header.Filename)
	fmt.Println(header.Header)

	defer file.Close()

	data := rw.NewData()
	data.Set("identifier", "local.jpg")

	testDir := &rw.File{
		Dir: "./test",
		Name: "testDir",
	}
	if err := testDir.Init(); err != nil {
		log.Fatal(err)
	}
	out, err := testDir.NewWriter(data)
	if err != nil {
		log.Fatal(err)
	}

	resize0 := &rw.Resize{
		Width:         500,
		Height:        0,
		Interpolation: "Lanczos3",
		Output:        "jpg",
		Name:          "resize",
	}
	resize0.Init()

	resize := &rw.Resize{
		Width:         300,
		Height:        0,
		Interpolation: "Lanczos3",
		Output:        "jpg",
		Name:          "resize",
	}
	resize.Init()

	zip := &rw.Shell{
		Cmd:  "gzip",
		Name: "zip",
	}
	zip.Init()

	unzip := &rw.Shell{
		Cmd:  "gzip -d",
		Name: "unzip",
	}
	unzip.Init()

	aes := &rw.AES256{
		Base64String: "ETl5QyPnHfi+vF4HrZfFvO2Julv4LVL7HNB1N7vkLGU=",
	}
	if err := aes.Init(); err != nil {
		log.Fatal(err)
	}

	p, err := rw.NewPipeline(
		[]rw.Encoder{resize0, resize, zip, unzip, zip, unzip, zip, unzip, zip, unzip, aes},
		file,
		data)
	if err != nil {
		fmt.Println(err)
	}
	err = p.Exec(out)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

func main() {
//	parseConfig()
	// GenCypher("aes256")

	http.HandleFunc("/", uploadHandler2)
	fmt.Println("server started")
	http.ListenAndServe(":6666", nil)
}
