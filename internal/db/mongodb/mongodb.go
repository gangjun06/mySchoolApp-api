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
var Category *mongo.Collection
var Post *mongo.Collection
var Calendar *mongo.Collection
var Schedule *mongo.Collection

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
	Category = db.Collection("category")
	Post = db.Collection("post")
	Calendar = db.Collection("calendar")
	Schedule = db.Collection("schedule")

	CheckErr(User.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"phone": 1},
			Options: options.Index().SetUnique(true),
		},
	))

	CheckErr(Calendar.Indexes().CreateMany(
		context.Background(),
		[]mongo.IndexModel{
			{Keys: bson.M{"year": 1}},
			{Keys: bson.M{"month": 1}},
			{Keys: bson.M{"day": 1}},
		},
	))

	CheckErr(Category.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
	))
}

func CheckErr(_dummy interface{}, err error) {
	if err != nil {
		log.Fatal(err)
	}
}
