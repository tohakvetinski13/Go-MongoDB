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

	quickstartDatabase := client.Database("quickstart")
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	episodesCollection := quickstartDatabase.Collection("episodes")

	// Reading All Documents from a Collection
	cursor, err := episodesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	// var episodes []bson.M
	// if err = cursor.All(ctx, &episodes); err != nil {
	// 	log.Fatal(err)
	// }
	// for _, episode := range episodes {
	// 	fmt.Println(episode["title"])
	// }
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var episode bson.M
		if err = cursor.Decode(&episode); err != nil {
			log.Fatal(err)
		}
		fmt.Println(episode)
	}

	// Reading a Single Document from a Collection
	var podcast bson.M
	if err = podcastsCollection.FindOne(ctx, bson.M{}).Decode(&podcast); err != nil {
		log.Fatal(err)
	}
	fmt.Println(podcast)

	// Querying Documents from a Collection with a Filter
	filterCursor, err := episodesCollection.Find(ctx, bson.M{"duration": 25})
	if err != nil {
		log.Fatal(err)
	}
	var episodesFiltred []bson.M
	if err = filterCursor.All(ctx, &episodesFiltred); err != nil {
		log.Fatal(err)
	}
	fmt.Println(episodesFiltred)

	// Sorting Documents in a Query
	opts := options.Find()
	opts.SetSort(bson.D{{"duration", -1}})
	sortCursor, err := episodesCollection.Find(ctx, bson.D{{"duration", bson.D{{"$gt", 24}}}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	var episodesSorted []bson.M
	if err = sortCursor.All(ctx, &episodesSorted); err != nil {
		log.Fatal(err)
	}
	fmt.Println(episodesSorted)
}
