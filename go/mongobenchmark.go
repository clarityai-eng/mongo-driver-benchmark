package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	Title  string `bson:"title"`
	Author string `bson:"author"`
}

func main() {

	start := time.Now()

	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := client.Database("test").Collection("pato")

	// Drop collection to start from scratch
	collection.Drop(context.Background())

	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go insertMany(&wg)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Insert documents took %s", elapsed)
	fmt.Println()

	// Now, find documents
	start = time.Now()

	batch_size := int32(1000)
	cursor, err := collection.Find(context.Background(), bson.M{"author": "George Orwell"}, &options.FindOptions{BatchSize: &batch_size})
	if err != nil {
		log.Fatal(err)
	}

	documents_found := 0
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		documents_found++
		var document bson.M
		if err = cursor.Decode(&document); err != nil {
			log.Fatal(err)
		}
	}

	elapsed = time.Since(start)
	fmt.Printf("Find %d documents took %s", documents_found, elapsed)

}

func insertMany(wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println("Starting goroutine")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("test").Collection("pato")

	insertedDocs := 0

	for i := 0; i < 6250; i++ {
		_1984 := Book{
			Title:  "1984",
			Author: "George Orwell",
		}
		animalfarm := Book{
			Title:  "Animal Farm",
			Author: "George Orwell",
		}
		greatgatsby := Book{
			Title:  "The Great Gatsby",
			Author: "F. Scott Fitzgerald",
		}

		books := []interface{}{_1984}

		for j := 0; j < 33; j++ {
			books = append(books, _1984, animalfarm, greatgatsby)
		}

		res, err := collection.InsertMany(context.Background(), books)

		if err != nil {
			log.Fatal(err)
		}
		insertedDocs += len(res.InsertedIDs)

	}

}
