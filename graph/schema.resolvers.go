package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Y-bro/go-graphql/graph/model"
	"github.com/Y-bro/go-graphql/internal/auth"
	"github.com/Y-bro/go-graphql/internal/links"
	"github.com/Y-bro/go-graphql/internal/users"
	"github.com/Y-bro/go-graphql/pkg/jwt"

	"github.com/Y-bro/go-graphql/graph/generated"
)

// CreateLink is the resolver for the createLink field.
func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {

	user := auth.ForContext(ctx)

	fmt.Print(user)

	if user == nil {
		return &model.Link{}, errors.New("login please")
	}

	link := links.Link{}
	link.Address = input.Address
	link.Title = input.Title
	link.User = user

	linkID := link.Save()

	graphQlUser := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}

	return &model.Link{ID: strconv.FormatInt(linkID, 10), Title: link.Address, Address: link.Address, User: graphQlUser}, nil

}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	user := users.User{}

	user.Username = input.Username
	user.Password = input.Password

	user.Create()

	token, err := jwt.GenerateToken(user.Username)

	if err != nil {
		return "", err
	}

	return token, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	user := users.User{}

	user.Password = input.Password
	user.Username = input.Username

	isAllowed := user.Authenticate()

	if !isAllowed {
		return "", errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.Username)

	if err != nil {
		return "", err
	}

	return token, nil
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)

	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(username)

	if err != nil {
		return "", err
	}

	return token, nil
}

// Links is the resolver for the links field.
func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	resLinks := []*model.Link{}

	dbLinks := links.GetAll()

	for _, link := range dbLinks {
		graphqlUser := model.User{
			ID:   link.User.ID,
			Name: link.User.Username,
		}
		resLinks = append(resLinks, &model.Link{ID: link.ID, Title: link.Title, Address: link.Address, User: &graphqlUser})

	}

	return resLinks, nil

}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
