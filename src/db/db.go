package db

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	MongoDB *mongo.Database
}

// DBConnection creates a new MongoDB client and connects to the server using the
// MONGODB environment variable. It returns a pointer to a DB struct, which
// contains the connected client and database.
func DBConnection() (*DB, error) {
	godotenv.Load()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// Set the URI for the MongoDB server from the MONGODB environment variable
	opts := options.Client().ApplyURI(os.Getenv("MONGODB")).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	// Connect to the "gdg-dev" database
	mdb := client.Database("gdg-dev")

	return &DB{
		MongoDB: mdb,
	}, nil
}
