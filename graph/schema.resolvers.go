package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/muesli/cache2go"
	"github.com/osang-school/backend/graph/generated"
	"github.com/osang-school/backend/graph/model"
	"github.com/osang-school/backend/graph/myerr"
	"github.com/osang-school/backend/internal/neis"
	"github.com/osang-school/backend/internal/post"
	"github.com/osang-school/backend/internal/session"
	"github.com/osang-school/backend/internal/user"
	"github.com/osang-school/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) SignIn(ctx context.Context, phone model.Phone, password string) (*model.ProfileWithToken, error) {
	userData, err := user.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	token, err := session.CreateToken(userData.ID, userData.Role, userData.Permissions)
	if err != nil {
		return nil, err
	}
	return &model.ProfileWithToken{
		Profile: user.UserToGqlType(userData),
		Token:   token,
	}, nil
}

func (r *mutationResolver) SignOut(ctx context.Context) (string, error) {
	if err := ctx.Value("data").(*session.Data).Expiry(); err != nil {
		return "", myerr.New(myerr.ErrServer, err.Error())
	}
	return "", nil
}

func (r *mutationResolver) VerifyPhone(ctx context.Context, number model.Phone) (string, error) {
	exits, err := user.CheckUserDup(string(number))
	if err != nil {
		return "", err
	} else if exits {
		return "", myerr.New(myerr.ErrDuplicate, "")
	}
	return "", user.PhoneVerifyCode(ctx.Value("ip").(string), string(number))
}

func (r *mutationResolver) CheckVerifyPhoneCode(ctx context.Context, number model.Phone, code string) (string, error) {
	return user.PhoneVerifyCheck(string(number), code)
}

func (r *mutationResolver) SetProfile(ctx context.Context, student *model.StudentProfileInput, teacher *model.TeacherProfileInput, officals *model.OfficalsProfileInput) (string, error) {
	randomStr := utils.CreateRandomString(6)
	cache := cache2go.Cache("profile")
	if student != nil {
		exits, err := user.CheckStudentDup(student.Grade, student.Class, student.Number)
		if err != nil {
			return "", err
		} else if exits {
			return "", myerr.New(myerr.ErrDuplicate, "")
		}
		cache.Add(randomStr, time.Minute*1, student)
	} else if teacher != nil {
		cache.Add(randomStr, time.Minute*1, teacher)
	} else if officals != nil {
		cache.Add(randomStr, time.Minute*1, officals)
	} else {
		return "", myerr.New(myerr.ErrBadRequest, "")
	}
	return randomStr, nil
}

func (r *mutationResolver) SignUp(ctx context.Context, input model.SignUpInput) (*model.ProfileWithToken, error) {
	phone, err := user.PhoneSignUpCheck(input.Phone)
	if err != nil {
		return nil, err
	}
	cache := cache2go.Cache("profile")
	res, err := cache.Value(input.Detail)
	if err != nil {
		return nil, myerr.New(myerr.ErrBadRequest, "")
	}
	detailData := res.Data()

	newUser := &user.User{
		Name:     input.Name,
		Phone:    phone,
		Status:   user.StatusWait,
		Password: utils.HashAndSalt(input.Password),
	}
	if input.Nickname != nil {
		newUser.Nickname = *input.Nickname
	}
	var resultDetail interface{}
	switch v := detailData.(type) {
	case *model.StudentProfileInput:
		newUser.Role = user.RoleStudent
		newUser.Student = &user.Student{
			Grade:  v.Grade,
			Class:  v.Class,
			Number: v.Number,
		}
		resultDetail = newUser.Student
	case *model.TeacherProfileInput:
		newUser.Role = user.RoleTeacher
		newUser.Teacher = &user.Teacher{
			Subject: v.Subject,
		}
		resultDetail = newUser.Teacher
	case *model.OfficalsProfileInput:
		newUser.Role = user.RoleOfficals
		newUser.Officals = &user.Officals{
			Role: v.Role,
		}
		if v.Description != nil {
			newUser.Officals.Description = *v.Description
		}
		resultDetail = newUser.Officals
	}

	id, err := user.SignUp(newUser)
	if err != nil {
		return nil, fmt.Errorf("Error While Signup")
	}

	phoneResult := model.Phone(newUser.Phone)
	profile := &model.Profile{
		ID:       model.ObjectID(id),
		Name:     newUser.Name,
		Nickname: newUser.Nickname,
		Phone:    &phoneResult,
		Status:   user.StatusToEnum(user.StatusWait),
		Detail:   user.DetailToUnion(resultDetail),
	}

	result := &model.ProfileWithToken{
		Profile: profile,
		Token:   "",
	}
	return result, nil
}

func (r *mutationResolver) CreateCategory(ctx context.Context, input model.NewCategory) (model.ObjectID, error) {
	convert := func(i model.UserRole) user.Role {
		switch i {
		case model.UserRoleStudent:
			return user.RoleStudent
		case model.UserRoleTeacher:
			return user.RoleTeacher
		case model.UserRoleOfficals:
			return user.RoleOfficals
		}
		return user.RoleOfficals
	}
	category := &post.Category{
		Name:          input.Name,
		ReqPermission: input.ReqPermission,
		AnonAble:      input.AnonAble,
	}
	for _, v := range input.ReadAbleRole {
		category.ReadAbleRole = append(category.ReadAbleRole, convert(v))
	}
	for _, v := range input.WriteAbleRole {
		category.WriteAbleRole = append(category.WriteAbleRole, convert(v))
	}

	id, err := post.NewCategory(category)
	if err != nil {
		return model.ObjectID(primitive.NilObjectID), err
	}
	return model.ObjectID(id), nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (model.ObjectID, error) {
	user := ctx.Value("data").(*session.Data)
	if category, err := post.GetCategory(primitive.ObjectID(input.Category)); err != nil {
		return model.ObjectID(primitive.NilObjectID), err
	} else {
		if ok := post.CheckUserPermission("write", category, user.Role, user.Permission); !ok {
			return model.ObjectID(primitive.NilObjectID), myerr.New(myerr.ErrPermission, "")
		}
	}

	id, err := post.NewPost(primitive.ObjectID(input.Category), user.ID, input.Title, input.Content)
	return model.ObjectID(id), err
}

func (r *mutationResolver) LikePost(ctx context.Context, input model.LikePostInput) (*string, error) {
	user := ctx.Value("data").(*session.Data)
	if category, err := post.GetCategoryByPost(primitive.ObjectID(input.Post)); err != nil {
		return nil, err
	} else {
		if ok := post.CheckUserPermission("read", category, user.Role, user.Permission); !ok {
			return nil, myerr.New(myerr.ErrPermission, "")
		}
	}
	data := ctx.Value("data").(*session.Data)
	err := post.PostLike(primitive.ObjectID(input.Post), data.ID, input.Status)
	return nil, err
}

func (r *mutationResolver) AddComment(ctx context.Context, input model.NewComment) (model.ObjectID, error) {
	user := ctx.Value("data").(*session.Data)
	if category, err := post.GetCategoryByPost(primitive.ObjectID(input.Post)); err != nil {
		return model.ObjectID(primitive.NilObjectID), err
	} else {
		if ok := post.CheckUserPermission("read", category, user.Role, user.Permission); !ok {
			return model.ObjectID(primitive.NilObjectID), myerr.New(myerr.ErrPermission, "")
		}
	}
	id, err := post.NewComment(primitive.ObjectID(input.Post), user.ID, input.Content)
	return model.ObjectID(id), err
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id model.ObjectID) (model.ObjectID, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) MyProfile(ctx context.Context) (*model.Profile, error) {
	userData := ctx.Value("user").(*user.User)
	return user.UserToGqlType(userData), nil
}

func (r *queryResolver) Cafeteria(ctx context.Context, filter *model.CafeteriaFilter) ([]*model.Cafeteria, error) {
	if filter == nil {
		filter = &model.CafeteriaFilter{}
	}

	return neis.GetCafeteria(filter)
}

func (r *queryResolver) Post(ctx context.Context, id model.ObjectID, comment *model.CommentFilter) (*model.Post, error) {
	userData := ctx.Value("data").(*session.Data)
	if category, err := post.GetCategoryByPost(primitive.ObjectID(id)); err != nil {
		return nil, err
	} else {
		if ok := post.CheckUserPermission("write", category, userData.Role, userData.Permission); !ok {
			return nil, myerr.New(myerr.ErrPermission, "")
		}
	}

	loadPost := true
	offset := 0
	limit := 20
	if comment != nil {
		if *comment.LoadOnlyComment {
			loadPost = false
		}
		if comment.Offset != nil {
			offset = *comment.Offset
		}
		if comment.Limit != nil {
			limit = *comment.Limit
		}
	}

	data, err := post.GetPost(primitive.ObjectID(id), loadPost, offset, limit)
	if err != nil {
		return nil, err
	}
	resultComment := []*model.Comment{}
	for _, v := range data.Comment {
		resultComment = append(resultComment, &model.Comment{
			ID:       model.ObjectID(v.ID),
			Content:  v.Content,
			CreateAt: model.Timestamp(v.CreateAt),
			UpdateAt: model.Timestamp(v.UpdateAt),
		})
	}
	return &model.Post{
		ID: model.ObjectID(data.ID),
		Category: &model.Category{
			ID:   model.ObjectID(data.CategoryData.ID),
			Name: data.CategoryData.Name,
		},
		Like:     0,
		IsLike:   true,
		Author:   user.UserToGqlType(data.AuthorData),
		Title:    data.Title,
		Content:  data.Content,
		CreateAt: model.Timestamp(data.CreateAt),
		UpdateAt: model.Timestamp(data.UpdateAt),
		Comment:  resultComment,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
