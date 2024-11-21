package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dnday/gdgoc-project/src/db"
	"github.com/dnday/gdgoc-project/src/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Book struct {
	ID           int       `json:"id" bson:"id"`
	Title        string    `json:"title" bson:"title"`
	Author       string    `json:"author" bson:"author"`
	Published_at string    `json:"published_at" bson:"published_at"`
	Updated_at   time.Time `json:"updated_at" bson:"updated_at"`
	Created_at   time.Time `json:"created_at" bson:"created_at"`
}

type BookRequest struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	Published_at string `json:"published_at"`
}
type BookResponse struct {
	Data []*Book `json:"data"`
}

func GetAllBook() (*BookResponse, error) {
	db, err := db.DBConnection()
	if err != nil {
		log.Default().Println(err.Error())
		return nil, errors.New("internal server error")
	}
	defer db.MongoDB.Client().Disconnect(context.TODO())

	coll := db.MongoDB.Collection("books")
	cur, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Default().Println(err.Error())
		return nil, errors.New("internal server error")
	}
	var books []*Book

	for cur.Next(context.TODO()) {
		var getbooks model.Book
		cur.Decode(&getbooks)
		books = append(books, &Book{
			ID:           getbooks.ID,
			Title:        getbooks.Title,
			Author:       getbooks.Author,
			Published_at: getbooks.Published_at,
			Created_at:   getbooks.Created_at,
			Updated_at:   getbooks.Updated_at,
		})
	}
	return &BookResponse{
		Data: books,
	}, nil
}

func AddBook(req io.Reader) error {
	var bookReq BookRequest
	err := json.NewDecoder(req).Decode(&bookReq)
	if err != nil {
		return errors.New("bad request")
	}

	db, err := db.DBConnection()
	if err != nil {
		log.Default().Println(err.Error())
		return errors.New("internal server error")
	}
	defer db.MongoDB.Client().Disconnect(context.TODO())

	coll := db.MongoDB.Collection("books")
	// Get the count of documents in the collection
	count, err := coll.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Default().Println(err.Error())
		return errors.New("internal server error")
	}

	// Insert the new book with an incremented ID
	_, err = coll.InsertOne(context.TODO(), model.Book{
		ID:           int(count) + 1,
		Title:        bookReq.Title,
		Author:       bookReq.Author,
		Published_at: bookReq.Published_at,
		Updated_at:   time.Now(),
		Created_at:   time.Now(),
	})
	if err != nil {
		log.Default().Println(err.Error())
		return errors.New("internal server error")
	}

	return nil

}

// GetBookByID returns a book by ID.
// If the book is not found, it sets the status code to 404 and returns an error.
func GetBookByID(w http.ResponseWriter, r *http.Request) (*Book, error) {
	// Get the book ID from the URL parameters.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		// If the id is not an integer, return an error.
		return nil, errors.New("id must be an integer")
	}

	// Connect to the database.
	db, err := db.DBConnection()
	if err != nil {
		// If there was an error connecting to the database, return the error.
		return nil, err
	}
	defer db.MongoDB.Client().Disconnect(context.TODO())

	// Get a collection of books.
	coll := db.MongoDB.Collection("books")

	// Create a filter to find the book by ID.
	filter := bson.D{{Key: "id", Value: id}}

	// Define a variable to hold the book.
	var book model.Book

	// Find the book by the filter and decode it into the variable.
	if err := coll.FindOne(context.Background(), filter).Decode(&book); err != nil {
		// If there was an error finding the book, handle it.
		if err == mongo.ErrNoDocuments {
			// If the book was not found, set the status code to 404 and return an error.
			w.WriteHeader(http.StatusNotFound)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Book not found",
			})
			return nil, nil
		}
		// If there was another type of error, return it.
		return nil, err
	}

	// Create a response variable to hold the book data.
	response := Book{
		ID:           book.ID,
		Title:        book.Title,
		Author:       book.Author,
		Published_at: book.Published_at,
		Created_at:   book.Created_at,
		Updated_at:   book.Updated_at,
	}

	// Return the response.
	return &response, nil
}
