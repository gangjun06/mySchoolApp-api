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
	"github.com/osang-school/backend/internal/session"
	"github.com/osang-school/backend/internal/user"
	"github.com/osang-school/backend/internal/utils"
)

func (r *mutationResolver) SignIn(ctx context.Context, phone model.Phone, password string) (*model.ProfileWithToken, error) {
	userData, err := user.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	token, err := session.CreateToken(userData.ID, &userData.Permissions)
	if err != nil {
		return nil, err
	}
	return &model.ProfileWithToken{
		Profile: user.UserToGqlType(userData),
		Token:   token,
	}, nil
}

func (r *mutationResolver) SignOut(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented"))
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

	profile := &model.Profile{
		ID:       model.ObjectID(id),
		Name:     newUser.Name,
		Nickname: newUser.Nickname,
		Phone:    model.Phone(newUser.Phone),
		Status:   user.StatusToEnum(user.StatusWait),
		Detail:   user.DetailToUnion(resultDetail),
	}

	result := &model.ProfileWithToken{
		Profile: profile,
		Token:   "",
	}
	return result, nil
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
