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
		// Handle GET request to retrieve all books
		data, err := service.GetAllBook()
		if err != nil {
			// Return internal server error if there's an issue retrieving books
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Set response content type to JSON
		w.Header().Add("Content-Type", "application/json")
		// Encode and send the list of books as JSON
		json.NewEncoder(w).Encode(data)
	case "POST":
		// Handle POST request to add a new book
		// Decode the JSON request body into a BookRequest
		resp, err := service.AddBook(r.Body)
		if err != nil {
			// Return internal server error if there's an issue adding the book
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Set response content type to JSON
		w.Header().Add("Content-Type", "application/json")
		// Set response status to Created
		w.WriteHeader(http.StatusCreated)
		// Encode and send the success message and added book data as JSON
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Book added successfully",
			"data":    resp,
		})
	case "PUT":
		// Handle PUT request to update an existing book
		err := service.UpdateBook(w, r)
		if err != nil {
			// Return internal server error if there's an issue updating the book
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Set response content type to JSON
		w.Header().Add("Content-Type", "application/json")
	case "DELETE":
		// Handle DELETE request to delete an existing book
		if err := service.DeleteBook(w, r); err != nil {
			// Return internal server error if there's an issue deleting the book
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Set response content type to JSON
		w.Header().Add("Content-Type", "application/json")
	default:
		// Handle unsupported HTTP methods
		log.Default().Println(http.StatusMethodNotAllowed)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// BookHandler handles HTTP requests for retrieving a single book by its ID.
// It supports the following HTTP method:
// - GET: Retrieves the details of a book by its ID.
// If an error occurs during the operation, an appropriate HTTP error response is returned.
func BookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		// Return method not allowed error if the request method is not GET
		log.Default().Println(http.StatusMethodNotAllowed)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	} else {
		// Handle GET request to retrieve a book by its ID
		book, err := service.GetBookByID(w, r)
		if err != nil {
			// Set response content type to JSON
			w.Header().Add("Content-Type", "application/json")
			// Encode and send the error message as JSON
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Book not found",
			})
			// Set response status to Not Found
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Set response content type to JSON
		w.Header().Add("Content-Type", "application/json")
		// Encode and send the book details as JSON
		if err := json.NewEncoder(w).Encode(book); err != nil {
			// Return internal server error if there's an issue encoding the book details
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
