package database

import (
	"ChatService/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func Connection(nameDB, nameCollection string) *mongo.Collection {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.ServiceConfig.MONGO_URL))
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database(nameDB).Collection(nameCollection)
	return collection
}
