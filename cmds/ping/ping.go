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
		fmt.Fprintf(os.Stderr, "Too few arguments")
		os.Exit(1)
	}
	listen = os.Args[1]
	pong = os.Args[2]
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, pong)
	})
	fmt.Fprintf(os.Stderr, "Starting\n")
	fmt.Printf("%s\n", http.ListenAndServe(listen, nil))
	fmt.Fprintf(os.Stderr, "Exiting\n")
}
