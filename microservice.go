package main

import (
	"fmt"
	"net/http"
	"os"

	"./api"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/api/echo", echo)

	http.HandleFunc("/api/books", api.BooksHandleFunction)
	http.HandleFunc("/api/books/", api.BookHandleFunction)

	http.ListenAndServe(port(), nil)
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "6060"
	}
	return ":" + port
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello Cloud Native Go Native.")
}

func echo(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query()["message"][0]

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, message)
}
