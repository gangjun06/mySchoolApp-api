package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/osang-school/backend/graph/generated"
	"github.com/osang-school/backend/graph/model"
	"github.com/osang-school/backend/internal/user"
)

func (r *mutationResolver) SignIn(ctx context.Context, phone string, password string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VerifyPhone(ctx context.Context, number model.Phone) (string, error) {
	err := user.PhoneVerifyCode("ip", string(number))
	return "", err
}

func (r *mutationResolver) SetProfile(ctx context.Context, student *model.StudentProfileInput, teacher *model.TeacherProfileInput, officals *model.OfficalsProfileInput) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SignUp(ctx context.Context, basic model.SignUpInput, detail string) (*model.Profile, error) {
	profile := &model.Profile{}
	profile.Name = "hello"
	return profile, nil
}

func (r *queryResolver) MyProfile(ctx context.Context) (*model.Profile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) IsValidPhoneVerifyCode(ctx context.Context, number model.Phone, code string) (bool, error) {
	correct, err := user.PhoneVerifyCheck(string(number), code)
	return correct, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
