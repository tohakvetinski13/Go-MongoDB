package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Podcast struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty"`
}

type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempy"`
	Podcast     primitive.ObjectID `bson:"podcast,omitempy"`
	Title       string             `bson:"title,omitempy"`
	Description string             `bson:"description,omitempy"`
	Duration    int32              `bson:"duration,omitempy"`
}

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

	database := client.Database("quickstart")
	podcastsCollection := database.Collection("podcasts")
	episodesCollection := database.Collection("episodes")

	var episodes []Episode
	cursor, err := episodesCollection.Find(ctx, bson.M{"duration": bson.D{{"$gt", 25}}})
	if err != nil {
		panic(err)
	}
	if err = cursor.All(ctx, &episodes); err != nil {
		panic(err)
	}
	fmt.Println(episodes)

	podcast := Podcast{
		Title:  "The Polyglot Developer",
		Author: "Nic Raboy",
		Tags:   []string{"development", "programming", "coding"},
	}
	insertResult, err := podcastsCollection.InsertOne(ctx, podcast)
	if err != nil {
		panic(err)
	}
	fmt.Println(insertResult.InsertedID)

}
