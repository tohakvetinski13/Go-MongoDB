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

// Podcast represent the schema for the "Podcasts" collection
type Podcast struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Author string             `bson:"author,omitempty"`
	Tags   []string           `bson:"tags,omitempty"`
}

// Episode represts the schema for the "Episodes" collection
type Episode struct {
	ID          primitive.ObjectID `bson:"_id,omitempy"`
	Podcast     primitive.ObjectID `bson:"podcast,omitempy"`
	Title       string             `bson:"title,omitempy"`
	Description string             `bson:"description,omitempy"`
	Duration    int32              `bson:"duration,omitempy"`
}

// PodcastEpisode represents an aggregation result-set for two collections
type PodcastEpisode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Podcast     Podcast            `bson:"podcast,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description,omitempty"`
	Duration    int32              `bson:"duration,omitempty"`
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
	episodesCollection := database.Collection("episodes")

	id, _ := primitive.ObjectIDFromHex("609c34b944d2b4a95dbecbf7")
	matchStage := bson.D{{"$match", bson.D{{"podcast", id}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$podcast"}, {"total", bson.D{{"$sum", "$duration"}}}}}}

	showInfoCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		panic(err)
	}
	var showsWithInfo []bson.M
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		panic(err)
	}
	fmt.Println(showsWithInfo)

	// lookupStage := bson.D{{"$lookup", bson.D{{"from", "podcasts"}, {"localField", "podcast"}, {"foreignField", "_id"}, {"as", "podcast"}}}}
	// unwindStage := bson.D{{"$unwind", bson.D{{"path", "$podcast"}, {"preserveNullAndEmptyArrays", false}}}}

	// showLoadedCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	// if err != nil {
	// 	panic(err)
	// }
	// var showsLoaded []bson.M
	// if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(showsLoaded)

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "podcasts"}, {"localField", "podcast"}, {"foreignField", "_id"}, {"as", "podcast"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$podcast"}, {"preserveNullAndEmptyArrays", false}}}}

	showLoadedStructCursor, err := episodesCollection.Aggregate(ctx, mongo.Pipeline{lookupStage, unwindStage})
	if err != nil {
		panic(err)
	}
	var showsLoadedStruct []PodcastEpisode
	if err = showLoadedStructCursor.All(ctx, &showsLoadedStruct); err != nil {
		panic(err)
	}
	fmt.Println(showsLoadedStruct)
}
