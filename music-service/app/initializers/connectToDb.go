package initializers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DBClient *mongo.Client
)

func ConnectToDB() {
	var mongoURL = "mongodb://" + os.Getenv("DB_ADMIN_LOGIN") + ":" + os.Getenv("DB_ADMIN_PWD") + "@music-service-db:27017/" + os.Getenv("DB_NAME")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		fmt.Println("Cannot connect to db: " + mongoURL)
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	DBClient = client
}
