package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/osang-school/backend/graph/generated"
	"github.com/osang-school/backend/graph/model"
)

func (r *mutationResolver) SignIn(ctx context.Context, phone string, password string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) VerifyPhone(ctx context.Context, number string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetProfile(ctx context.Context, student *model.StudentProfileInput, teacher *model.TeacherProfileInput, officals *model.OfficalsProfileInput) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SignUp(ctx context.Context, basic model.SignUpInput, detail string) (model.Profile, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
