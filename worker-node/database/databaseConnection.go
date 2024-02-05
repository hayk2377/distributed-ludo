package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GameState struct {
	BsonID   primitive.ObjectID `bson:"_id"`
	GameCode *string            `json:"gameCode"`
	State    *string            `json:"state"`
}

func DBinstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URL")

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBinstance()
var COLLECTION_NAME = "gameState"

var gameStateCollection *mongo.Collection = Client.Database("cluster0").Collection(COLLECTION_NAME)

func GetGameState(gameCode string) (string, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var gameState GameState
	err := gameStateCollection.FindOne(ctx, bson.M{"gameCode": gameCode}).Decode(&gameState)
	defer cancel()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("no game state found for game code: %v", gameCode)
		}
		fmt.Println("error getting state", err)
		return "", err
	}
	return *gameState.State, nil
}

func UpsertGameState(gameCode string, state string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"state": state}}
	opts := options.FindOneAndUpdate().SetUpsert(true)

	var gameState GameState
	result := gameStateCollection.FindOneAndUpdate(ctx, bson.M{"gameCode": gameCode}, update, opts)
	err := result.Decode(&gameState)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println("Error setting state:", err)
		return fmt.Errorf("error updating game state: %v", err)
	}
	fmt.Println("succesfully upserted game state", gameCode, state)

	return nil
}
