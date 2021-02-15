package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/osang-school/backend/internal/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var db *mongo.Database

var User *mongo.Collection

func Init() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(conf.Get().MongoDB))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	db = client.Database("osang")
	User = db.Collection("user")

	_, err = User.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"phone": 1},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
}
