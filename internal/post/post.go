package post

import (
	"time"

	"github.com/osang-school/backend/graph/myerr"
	"github.com/osang-school/backend/internal/db/mongodb"
	"github.com/osang-school/backend/internal/user"
	"github.com/osang-school/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Category struct {
		ID                  primitive.ObjectID `bson:"_id,omitempty"`
		Name                string             `bson:"name,omitempty"`
		ReqPermission       []string           `bson:"reqPermission,omitempty"`
		ReqManagePermission []string           `bson:"reqManagePermission,omitempty"`
		AnonAble            bool               `bson:"anonAble,omitempty"`
		ReadAbleRole        []user.Role        `bson:"readAbleRole,omitempty"`
		WriteAbleRole       []user.Role        `bson:"writeAbleRole_id,omitempty"`
	}
	PostStatus uint8
	Post       struct {
		ID           primitive.ObjectID   `bson:"_id,omitempty"`
		Anon         bool                 `bson:"anon,omitempty"`
		Category     primitive.ObjectID   `bson:"category,omitempty"`
		CategoryData *Category            `bson:"categoryData,omitempty"`
		Like         []primitive.ObjectID `bson:"likeUsers,omitempty"`
		LikeCnt      int                  `bson:"likeCnt,omitempty"`
		IsLike       bool                 `bson:"isLike,omitempty"`
		Author       primitive.ObjectID   `bson:"author,omitempty"`
		AuthorData   *user.User           `bson:"authorData,omitempty"`
		Status       PostStatus           `bson:"status,omitempty"`
		Title        string               `bson:"title,omitempty"`
		Content      string               `bson:"content,omitempty"`
		CreateAt     time.Time            `bson:"createAt,omitempty"`
		UpdateAt     time.Time            `bson:"updateAt,omitempty"`
		Comment      []*Comment           `bson:"comment,omitempty"`
	}
	Comment struct {
		ID         primitive.ObjectID `bson:"_id,omitempty"`
		Author     primitive.ObjectID `bson:"author,omitempty"`
		AuthorData *user.User         `bson:"authorData,omitempty"`
		Status     PostStatus         `bson:"status,omitempty"`
		Content    string             `bson:"content,omitempty"`
		CreateAt   time.Time          `bson:"createAt,omitempty"`
		UpdateAt   time.Time          `bson:"updateAt,omitempty"`
	}
)

const (
	StatusNormal PostStatus = iota + 1
	StatusDeleted
	StatusReported
)

func CheckUserPermission(purpose string, category *Category, userRole user.Role, userPermission []string) bool {
	if ok := utils.ArrayHasItem(userPermission, "admin"); ok {
		return true
	}
	switch purpose {
	case "write":
		if len(category.ReqPermission) > 0 {
			if ok := utils.HasPermission(category.ReqPermission, userPermission); ok {
				return true
			}
			return false
		}
		if ok := func() bool {
			for _, v := range category.WriteAbleRole {
				if userRole == v {
					return true
				}
			}
			return false
		}(); !ok {
			return false
		}
	case "read":
		if ok := func() bool {
			for _, v := range category.ReadAbleRole {
				if userRole == v {
					return true
				}
			}
			return false
		}(); !ok {
			return false
		}
	case "manage":
		if len(category.ReqManagePermission) > 0 {
			if ok := utils.HasPermission(category.ReqPermission, userPermission); ok {
				return true
			}
			return false
		}
	}
	return false
}

func NewCategory(category *Category) (primitive.ObjectID, error) {
	result, err := mongodb.Category.InsertOne(nil, category)
	if err != nil {
		return primitive.NilObjectID, myerr.New(myerr.ErrServer, err.Error())
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func GetCategory(id primitive.ObjectID) (*Category, error) {
	filter := bson.M{"_id": id}
	var result Category
	err := mongodb.Category.FindOne(nil, filter).Decode(&result)
	if err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	return &result, nil
}

func GetCategoryByPost(id primitive.ObjectID) (*Category, error) {
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.M{"_id": id}}},
		bson.D{{"$lookup", bson.D{{"from", "category"}, {"localField", "category"}, {"foreignField", "_id"}, {"as", "category"}}}},
		bson.D{{"$unwind", "$category"}},
		bson.D{{"$project", bson.M{"category": 1}}},
	}
	cursor, err := mongodb.Post.Aggregate(nil, pipeline)
	if err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	var data []struct {
		Category Category `bson:"category"`
	}
	if err := cursor.All(nil, &data); err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	if len(data) < 1 {
		return nil, myerr.New(myerr.ErrNotFound, "post not found")
	}
	return &data[0].Category, nil
}

func CategoryExits(id primitive.ObjectID) (bool, error) {
	var result Category
	err := mongodb.Category.FindOne(nil, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(&result)
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewPost(categoryID, author primitive.ObjectID, title, content string, anon bool) (primitive.ObjectID, error) {
	post := &Post{
		Author:   author,
		Category: categoryID,
		Title:    title,
		Content:  content,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
		Anon:     anon,
		Status:   StatusNormal,
	}
	result, err := mongodb.Post.InsertOne(nil, post)
	if err != nil {
		return primitive.NilObjectID, myerr.New(myerr.ErrServer, err.Error())
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// GetPost (postid, loadPost, commentOffset, commentLimit)
func GetPost(id, userID primitive.ObjectID, loadPost bool, comment ...int) (*Post, error) {
	commentOffset := 0
	commentLimit := 0
	if len(comment) == 2 {
		commentOffset = comment[0]
		commentLimit = comment[1]
	}
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.M{"_id": id}}},
	}
	if loadPost {
		pipeline = append(pipeline, mongo.Pipeline{
			bson.D{{"$lookup", bson.D{{"from", "category"}, {"localField", "category"}, {"foreignField", "_id"}, {"as", "categoryData"}}}},
			bson.D{{"$lookup", bson.D{{"from", "user"}, {"localField", "author"}, {"foreignField", "_id"}, {"as", "authorData"}}}},
			bson.D{{"$unwind", "$categoryData"}},
			bson.D{{"$unwind", "$authorData"}},
		}...)
	}

	if commentLimit < 1 {
		if loadPost {
			pipeline = append(pipeline, bson.D{{"$project", bson.M{"comment": 0}}})
		}
	} else {
		pipeline = append(pipeline, mongo.Pipeline{
			bson.D{{"$unwind", bson.M{
				"path":                       "$comment",
				"preserveNullAndEmptyArrays": true,
			}}},
			bson.D{{"$skip", commentOffset}},
			bson.D{{"$limit", commentLimit}},
			bson.D{{"$lookup", bson.D{{"from", "user"}, {"localField", "comment.author"}, {"foreignField", "_id"}, {"as", "comment.authorData"}}}},
			bson.D{{"$unwind", bson.M{
				"path":                       "$comment.authorData",
				"preserveNullAndEmptyArrays": true,
			}}},
			bson.D{{"$group", bson.M{
				"_id":     "$_id",
				"root":    bson.M{"$mergeObjects": "$$ROOT"},
				"comment": bson.M{"$push": "$comment"},
			}}},
			bson.D{{"$replaceRoot", bson.M{
				"newRoot": bson.M{
					"$mergeObjects": bson.A{"$root", "$$ROOT"},
				},
			},
			}},
			bson.D{{"$project", bson.M{
				"root": 0,
			}}},
			bson.D{{"$set", bson.M{
				"likeCnt": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": bson.A{"$likeUsers", nil}},
						"then": bson.M{"$size": "$likeUsers"},
						"else": 0,
					},
				},
				"isLike": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": bson.A{"$likeUsers", nil}},
						"then": bson.M{"$in": bson.A{userID, "$likeUsers"}},
						"else": false,
					},
				},
			}}},
		}...)
	}
	pipeline = append(pipeline, bson.D{{"$project", bson.M{"like": 0}}})

	cursor, err := mongodb.Post.Aggregate(nil, pipeline)
	if err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	var post []Post
	if err := cursor.All(nil, &post); err != nil {
		return nil, myerr.New(myerr.ErrServer, err.Error())
	}
	if len(post) < 1 {
		return nil, myerr.New(myerr.ErrNotFound, "")
	}
	if len(post[0].Comment) == 1 && post[0].Comment[0].Content == "" {
		post[0].Comment = []*Comment{}
	}
	if post[0].Anon {
		post[0].AuthorData = &user.User{
			Name:   "익명",
			Role:   user.RoleAnon,
			Status: user.StatusUser,
			Phone:  "01000000000",
		}
	}
	return &post[0], nil
}

func PostLike(postID, user primitive.ObjectID, status bool) error {
	filter := bson.M{"_id": postID}
	update := bson.M{"$addToSet": bson.M{
		"likeUsers": user,
	}}
	if !status {
		update = bson.M{"$pull": bson.M{
			"likeUsers": user,
		}}
	}
	_, err := mongodb.Post.UpdateMany(nil, filter, update)
	if err != nil {
		return myerr.New(myerr.ErrServer, err.Error())
	}
	return nil
}

func NewComment(postID, author primitive.ObjectID, content string) (primitive.ObjectID, error) {
	filter := bson.M{"_id": postID}
	objectID := primitive.NewObjectID()
	update := bson.M{"$push": bson.M{
		"comment": bson.M{
			"_id":      objectID,
			"author":   author,
			"content":  content,
			"createAt": time.Now(),
			"updateAt": time.Now(),
		},
	}}
	_, err := mongodb.Post.UpdateOne(nil, filter, update)
	if err != nil {
		return primitive.NewObjectID(), myerr.New(myerr.ErrServer, err.Error())
	}
	return objectID, nil
}
