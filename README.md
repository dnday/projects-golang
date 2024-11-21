# gdgoc
This is a simple RESTful API written in Go language for storing and retrieving book information.

The program uses `gorilla/mux` as the router and `go.mongodb.org/mongo-driver` as the MongoDB driver.

## Installation
Ensure you have Go installed on your machine.

Clone the repository and run the following command in the terminal:

    go build

## Usage
Run the program with:

    ./src/main

You can use the API with the following endpoints:

    GET /api/books
    POST /api/books
    PUT /api/books/{id}
    DELETE /api/books/{id}
    GET /api/books/{id}

The API returns results in JSON format containing the following information:

    id: the ID of the book in MongoDB (not ObjectID)
    title: the title of the book
    author: the author of the book
    published_at: the publication date of the book
    updated_at: the last update date of the book data
    created_at: the creation date of the book data

## Environment variables
Set the following environment variable:

    MONGODB: the URI of your MongoDB server

For example, you can set it in your .bashrc file:

    export MONGODB="mongodb://localhost:27017"
