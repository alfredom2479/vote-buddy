package main

import (
	"fmt"
	"net/http"
)

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleGetRoot func is hit")

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "templates/index.html")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("testHandler function entered")
	fmt.Fprintf(w, "Does this work?")
}
