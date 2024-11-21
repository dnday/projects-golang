package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dnday/gdgoc-project/src/db"
	"github.com/dnday/gdgoc-project/src/logger"
	"github.com/dnday/gdgoc-project/src/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Book struct {
	ID           string    `json:"id" bson:"_id"`
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
			ID:           getbooks.ID.Hex(),
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
	_, err = coll.InsertOne(context.TODO(), model.Book{
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

func GetBookByID(w http.ResponseWriter, r *http.Request) (*Book, error) {
	param := mux.Vars(r)["id"]
	_id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		logger.ResponseLogger(r, http.StatusBadRequest, err.Error())
		http.Error(w, "server: bad request", http.StatusBadRequest)
		return nil, err
	}

	db, err := db.DBConnection()
	if err != nil {
		logger.ResponseLogger(r, http.StatusInternalServerError, err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return nil, err
	}
	defer db.MongoDB.Client().Disconnect(context.TODO())

	coll := db.MongoDB.Collection("books")
	filter := bson.D{{Key: "_id", Value: _id}}

	var book model.Book
	if err := coll.FindOne(context.Background(), filter).Decode(&book); err != nil {
		if err == mongo.ErrNoDocuments {
			logger.ResponseLogger(r, http.StatusNotFound, "book not found")
			http.Error(w, "Book not found", http.StatusNotFound)
			return nil, err
		}
		logger.ResponseLogger(r, http.StatusInternalServerError, err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return nil, err
	}

	response := Book{
		ID:           book.ID.Hex(),
		Title:        book.Title,
		Author:       book.Author,
		Published_at: book.Published_at,
		Created_at:   book.Created_at,
		Updated_at:   book.Updated_at,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.ResponseLogger(r, http.StatusInternalServerError, err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return nil, err
	}
	return &response, nil
}
