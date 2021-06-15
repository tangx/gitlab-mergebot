package confgitlab

import (
	"context"

	"github.com/go-courier/httptransport/httpx"
	pkgerrors "github.com/pkg/errors"
)

type GetUserInput struct {
	httpx.MethodGet
	User
}

func (GetUserInput) Path() string {
	return "/api/v4/users"
}

func (g *GitlabV4) GetUser(input GetUserInput) ([]*User, error) {
	user := []*User{}
	ctx := context.Background()
	_, err := g.c.Do(ctx, &input, g.metas).Into(&user)
	if err != nil {
		return nil, pkgerrors.Wrapf(err, "gitlab get user failed: %v", err)
	}

	return user, nil
}

func (g *GitlabV4) GetUserByName(name string) ([]*User, error) {
	input := GetUserInput{
		User: User{Username: name},
	}
	return g.GetUser(input)
}
