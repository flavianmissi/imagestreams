package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World! 123")
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
