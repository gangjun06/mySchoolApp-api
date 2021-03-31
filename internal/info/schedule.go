package info

import (
	"github.com/osang-school/backend/graph/myerr"
	"github.com/osang-school/backend/internal/db/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Schedule struct {
		ID          primitive.ObjectID `bson:"_id"`
		Subject     string             `bson:"subject,omitempty"`
		Teacher     string             `bson:"teacher,omitempty"`
		Description string             `bson:"description,omitempty"`
		ClassRoom   string             `bson:"classRoom,omitempty"`
		Grade       uint               `bson:"grade,omitempty"`
		Class       uint               `bson:"class,omitempty"`
		Dow         uint               `bson:"dow,omitempty"`
		Period      uint               `bson:"period,omitempty"`
	}
	UpdateScheduleInput struct {
		Grade, Class, Dow, Period                uint
		Subject, Teacher, Description, ClassRoom string
	}
)

func InsertSchedules(input []*Schedule) error {
	var newValue []interface{}

	for _, v := range input {
		newValue = append(newValue, v)
	}

	_, err := mongodb.Schedule.InsertMany(nil, newValue)
	if err != nil {
		return myerr.New(myerr.ErrServer, err.Error())
	}

	return nil
}

func UpdateSchedule(input *UpdateScheduleInput) error {
	filter := bson.M{
		"grade":  input.Grade,
		"class":  input.Class,
		"dow":    input.Dow,
		"period": input.Period,
	}
	update := bson.M{
		"$set": bson.M{
			"grade":       input.Grade,
			"class":       input.Class,
			"dow":         input.Dow,
			"period":      input.Period,
			"subject":     input.Subject,
			"teacher":     input.Teacher,
			"description": input.Description,
			"classroom":   input.ClassRoom,
		},
	}
	_, err := mongodb.Schedule.UpdateOne(nil, filter, &update, options.Update().SetUpsert(true))
	if err != nil {
		return myerr.New(myerr.ErrServer, err.Error())
	}

	return nil
}

func FindSchedule(grade, class, dow uint) ([]*Schedule, error) {
	filter := bson.M{
		"grade": grade,
		"class": class,
		"dow":   dow,
	}
	cursor, err := mongodb.Schedule.Find(nil, filter)
	if err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	var data []*Schedule
	if err := cursor.All(nil, &data); err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	return data, nil
}

func DeleteSchedule(grade, class, dow, period uint) error {
	filter := bson.M{
		"grade":  grade,
		"class":  class,
		"dow":    dow,
		"period": period,
	}
	_, err := mongodb.Schedule.DeleteOne(nil, filter)
	return err
}
