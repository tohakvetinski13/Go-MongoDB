package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Creating Connection with MongoDB
	// ATLAS_URI your personal URI in Atlas Cluster
	client, err := mongo.NewClient(options.Client().ApplyURI(ATLAS_URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	defer client.Disconnect(ctx)

	database := client.Database("quickstart")
	podcastsCollection := database.Collection("podcasts")
	episodesCollection := database.Collection("episodes")

	// Deleting a Single Document from a MongoDB Collection
	result, err := podcastsCollection.DeleteOne(
		ctx,
		bson.M{"title": "The Polyglot Developer Podcast"},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)

	// Deleting Many Documents from a MongoDB Collection
	result, err = episodesCollection.DeleteMany(
		ctx,
		bson.M{"duration": 25},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)

	// Dropping a MongoDB Collection and All Documents within the Collection
	if err = podcastsCollection.Drop(ctx); err != nil {
		log.Fatal(err)
	}

	if err = episodesCollection.Drop(ctx); err != nil {
		log.Fatal(err)
	}
}
