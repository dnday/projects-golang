package main

import (
	"log"
	"net/http"

	api "github.com/dnday/gdgoc-project/api"
	"github.com/gorilla/mux"
)

func main() {
	h := mux.NewRouter()
	api.BooksRoute("/api/books", h)
	log.Default()
	h.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API!"))
	})
	s := &http.Server{
		Addr:    ":8000",
		Handler: h,
	}

	log.Default().Println("Server is running on port 8000")
	log.Fatal(s.ListenAndServe())

}
