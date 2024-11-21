package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dnday/gdgoc-project/src/service"
)

func ListBooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data, err := service.GetAllBook()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
		return
	case "POST":
		// Close the request body when the function returns
		defer r.Body.Close()

		// Decode the JSON request body into a map
		var book map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Marshal the book map into JSON bytes
		bookBytes, err := json.Marshal(book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Call the service to add the book
		err = service.AddBook(bytes.NewReader(bookBytes))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set the response header to JSON and return a success message
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(
			map[string]interface{}{
				"message": "Book created successfully",
				"data":    book,
			},
		)
		return
	case "PUT":

	default:
		log.Default().Println(http.StatusMethodNotAllowed)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func BookHandler(w http.ResponseWriter, r *http.Request) {
	book, err := service.GetBookByID(w, r)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
