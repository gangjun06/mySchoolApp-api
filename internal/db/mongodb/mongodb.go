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

var mClient *mongo.Client
var mdb *mongo.Database

var mUser *mongo.Collection

func Init() {
	client, err := mongo.NewClient(options.Client().ApplyURI(conf.Get().MongoDB))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	mClient = client
	mdb = client.Database("osang")
	mUser = mdb.Collection("user")

	_, err = mUser.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
}
