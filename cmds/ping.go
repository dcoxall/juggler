package main

import (
	"fmt"
	"net/http"
	"os"
)

var (
	// This is the configurable response to every request
	pong string
	// This is the address to listen on
	listen string
)

func init() {
	// Check that we have the right number of arguments
	if len(os.Args) != 3 {
		os.Exit(1)
	}
	listen = os.Args[1]
	pong = os.Args[2]
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, pong)
	})
	http.ListenAndServe(listen, nil)
}
