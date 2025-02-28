package main

import (
	"fmt"
	"log"
	"net/http"
	"webFramework/waku"
)

func main() {
	r := waku.NewEngine()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	log.Fatal(r.Run(":8000"))

}
