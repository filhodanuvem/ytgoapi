package post

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/filhodanuvem/ytgoapi/internal"
)

var ErrPostBodyEmpty = errors.New("post body is empty")
var ErrPostBodyExceedsLimit = errors.New("post body exceeds limit")
var ErrPostNotFound = errors.New("post not found")
var ErrPostUsernameEmpty = errors.New("post username empty")
var ErrIdEmpty = errors.New("id empty")

type Service struct {
	Repository Repository
}

func (p Service) Create(ctx context.Context, post internal.Post) (internal.Post, error) {
	if post.Body == "" {
		return internal.Post{}, ErrPostBodyEmpty
	}

	if post.Username == "" {
		return internal.Post{}, ErrPostUsernameEmpty
	}

	if utf8.RuneCountInString(post.Body) > 140 {
		return internal.Post{}, ErrPostBodyExceedsLimit
	}

	return p.Repository.Insert(ctx, post)
}

func (s Service) Delete(ctx context.Context, id string) error {
	return s.Repository.Delete(ctx, id)
}

func (s Service) FindOneByID(ctx context.Context, id string) (internal.Post, error) {
	return s.Repository.FindOneByID(ctx, id)
}

func (s Service) FindAll(ctx context.Context) ([]internal.Post, error) {
	posts, err := s.Repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s Service) Update(ctx context.Context, params *internal.ParamsUpdatePost) error {
	post, err := params.Validate()
	if err != nil {
		return err
	}

	err = s.Repository.Update(ctx, post)

	if err != nil {
		return err
	}

	return nil
