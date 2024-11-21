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

// main is the entry point of the program.
func main() {
	// Create a new Router.
	h := mux.NewRouter()

	// Add the BooksRoute to the router.
	BooksRoute("/api/books", h)

	// Add a route for the root url.
	h.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// Just return a simple message.
		w.Write([]byte("Hello from API!"))
	})

	// Create a new Server.
	s := &http.Server{
		Addr:    ":8000", // Listen on port 8000.
		Handler: h,
	}

	// Log a message when the server starts.
	log.Default().Println("Server is running on port 8000")

	// Start the server.
	// This will block until the server is stopped.
	log.Fatal(s.ListenAndServe())
}
