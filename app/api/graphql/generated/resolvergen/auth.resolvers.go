package resolvergen

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.27

import (
	"context"
	"fmt"
	"gogql/app/api/graphql/generated/graph"
	"gogql/app/models"
)

// GenerateOtp is the resolver for the generateOTP field.
func (r *mutationResolver) GenerateOtp(ctx context.Context, input *graph.OTPRequest) (*string, error) {
	panic(fmt.Errorf("not implemented: GenerateOtp - generateOTP"))
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input graph.LoginRequest) (*models.Auther, error) {
	panic(fmt.Errorf("not implemented: Login - login"))
}

// Auther is the resolver for the auther field.
func (r *queryResolver) Auther(ctx context.Context) (*models.Auther, error) {
	panic(fmt.Errorf("not implemented: Auther - auther"))
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }