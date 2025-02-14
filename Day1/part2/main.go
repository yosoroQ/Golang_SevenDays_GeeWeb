package main

import (
	"fmt"
	"log"
	"net/http"
	// "os"
)

type Engine struct {
}

func main() {
	enginge := new(Engine)
	log.Fatal(http.ListenAndServe(":8080", enginge))
}

func (enginge *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
