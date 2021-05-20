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
	collection := client.Database("test").Collection("go")

	// Drop collection to start from scratch
	collection.Drop(context.Background())

	// Insert documents using 8 goroutines
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
	cursor, err := collection.Find(context.Background(), bson.M{"author": "Isaac Asimov"}, &options.FindOptions{BatchSize: &batch_size})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var document bson.M
		if err = cursor.Decode(&document); err != nil {
			log.Fatal(err)
		}

		var _ = document["title"]
		var _ = document["author"]
	}

	elapsed = time.Since(start)
	fmt.Printf("Find documents took %s", elapsed)

}

func insertMany(wg *sync.WaitGroup) {

	defer wg.Done()

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("test").Collection("go")

	dune := Book{
		Title:  "Dune",
		Author: "Frank Herbert",
	}
	i_robot := Book{
		Title:  "I, Robot",
		Author: "Isaac Asimov",
	}
	foundation := Book{
		Title:  "Foundation",
		Author: "Isaac Asimov",
	}
	brave_new_world := Book{
		Title:  "Brave New World",
		Author: "Aldous Huxley",
	}

	books := []interface{}{dune}

	for j := 0; j < 33; j++ {
		books = append(books, i_robot, foundation, brave_new_world)
	}

	for i := 0; i < 6250; i++ {

		_, err := collection.InsertMany(context.Background(), books)

		if err != nil {
			log.Fatal(err)
		}

	}

}
