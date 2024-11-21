package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dnday/gdgoc-project/src/service"
)

// ListBooksHandler handles HTTP requests for listing, adding, updating, and deleting books.
// It supports the following HTTP methods:
// - GET: Retrieves a list of all books.
// - POST: Adds a new book. The request body should contain the book details in JSON format.
// - PUT: Updates an existing book. The request body should contain the updated book details in JSON format.
// - DELETE: Deletes an existing book.
// If an error occurs during any operation, an appropriate HTTP error response is returned.
func ListBooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data, err := service.GetAllBook()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	case "POST":
		// Decode the JSON request body into a BookRequest
		resp, err := service.AddBook(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Book added successfully",
			"data":    resp,
		})
	case "PUT":
		err := service.UpdateBook(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")

	case "DELETE":
		if err := service.DeleteBook(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
	default:
		log.Default().Println(http.StatusMethodNotAllowed)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func BookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Default().Println(http.StatusMethodNotAllowed)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	} else {
		book, err := service.GetBookByID(w, r)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Book not found",
			})
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
