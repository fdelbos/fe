//
// server.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  8 2014.
//

package main

import (
	//"github.com/fdelbos/fe/rw"
	//"net/http"
	"fmt"
	"os"
)

func main() {
	_, err := parseConfig()
	if err != nil {
		fmt.Println("Configuration error:")
		fmt.Println(err)
		os.Exit(-1)
	}

	// http.HandleFunc("/", uploadHandler2)
	// fmt.Println("server started")
	// http.ListenAndServe(":6666", nil)

}
