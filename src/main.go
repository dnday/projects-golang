package main

import (
	"log"
	"net/http"

	api "github.com/dnday/gdgoc-project/api"
	"github.com/gorilla/mux"
)

func BooksRoute(prefix string, r *mux.Router) {
	b := r.PathPrefix(prefix).Subrouter()
	b.HandleFunc("", api.ListBooksHandler).Methods("GET", "POST")
	b.HandleFunc("/{id}", api.ListBooksHandler).Methods("PUT", "DELETE")
	b.HandleFunc("/{id}", api.BookHandler).Methods("GET")
}
func main() {
	h := mux.NewRouter()

	BooksRoute("/api/books", h)
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
