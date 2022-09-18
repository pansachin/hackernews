package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"strconv"

	"github.com/pansachin/hackernews/graph/generated"
	"github.com/pansachin/hackernews/graph/model"
	"github.com/pansachin/hackernews/internal/auth"
	"github.com/pansachin/hackernews/internal/links"
	"github.com/pansachin/hackernews/internal/users"
	"github.com/pansachin/hackernews/pkg/jwt"
)

// CreateLink is the resolver for the createLink field.
func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	// panic(fmt.Errorf("not implemented: CreateLink - createLink"))
	u := auth.ForContext(ctx)

	var link links.Link
	link.Address = input.Address
	link.Title = input.Title
	link.User = u
	linkID := link.Save()

	gqUser := &model.User{
		Name: u.UserName,
		ID:   u.ID,
	}
	return &model.Link{
		ID:      strconv.FormatInt(linkID, 10),
		Title:   link.Title,
		Address: link.Address,
		User:    gqUser,
	}, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	// panic(fmt.Errorf("not implemented: CreateUser - createUser"))
	var user users.User

	user.UserName = input.Name
	user.Password = input.Password

	user.Create()
	token, err := jwt.GenerateToken(input.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	// panic(fmt.Errorf("not implemented: Login - login"))
	u := users.User{
		UserName: input.Username,
		Password: input.Password,
	}
	ok := u.Authenticate()

	if !ok {
		return "", &users.WrongUserOrPasswordError{}
	}
	token, err := jwt.GenerateToken(input.Username)
	if err != nil {
		return "", err
	}

	return token, err
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, input string) (string, error) {
	// panic(fmt.Errorf("not implemented: RefreshToken - refreshToken"))
	username, err := jwt.ParseToken(input)
	if err != nil {
		return "", errors.New("access denied")
	}

	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}

	return token, err

}

// Links is the resolver for the links field.
func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	// panic(fmt.Errorf("not implemented: Links - links"))
	// var links []*model.Link
	// dummyLink := model.Link{
	// 	Title:   "our dummy link",
	// 	Address: "https://address.org",
	// 	User:    &model.User{Name: "admin"},
	// }
	// links = append(links, &dummyLink)

	// return links, nil

	var resultLink []*model.Link

	dbLinks := links.GetAll()

	for _, link := range dbLinks {
		resultLink = append(resultLink, &model.Link{
			ID:      link.ID,
			Title:   link.Title,
			Address: link.Address,
			User: &model.User{
				Name: link.User.UserName,
				ID:   link.User.ID,
			},
		})
	}

	return resultLink, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
