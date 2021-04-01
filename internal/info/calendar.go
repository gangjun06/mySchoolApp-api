package info

import (
	"github.com/osang-school/backend/graph/myerr"
	"github.com/osang-school/backend/internal/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Calendar struct {
		ID          primitive.ObjectID `bson:"_id"`
		Year        uint               `bson:"year"`
		Month       uint               `bson:"month"`
		Day         uint               `bson:"day"`
		Title       string             `bson:"title"`
		Description string             `bson:"description,omitempty"`
		Icon        string             `bson:"icon,omitempty"`
	}
)

func NewCalendar(year, month, day uint, title, description, icon string) (primitive.ObjectID, error) {
	newItem := Calendar{
		ID:          primitive.NewObjectID(),
		Title:       title,
		Description: description,
		Icon:        icon,
		Year:        year,
		Month:       month,
		Day:         day,
	}
	result, err := mongodb.Calendar.InsertOne(nil, &newItem)
	if err != nil {
		return primitive.NilObjectID, myerr.New(myerr.ErrServer, err.Error())
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func FindCalendar(year, month uint) ([]*Calendar, error) {
	filter := bson.M{
		"year":  year,
		"month": month,
	}
	cursor, err := mongodb.Calendar.Find(nil, filter, options.Find().SetSort(bson.M{
		"day": 1,
	}))
	if err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	var data []*Calendar
	if err := cursor.All(nil, &data); err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	return data, nil
}

func DeleteCalendar(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := mongodb.Calendar.DeleteOne(nil, filter)
	return err
}
