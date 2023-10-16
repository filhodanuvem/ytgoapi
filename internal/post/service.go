package post

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/filhodanuvem/ytgoapi/internal"
	"github.com/google/uuid"
)

var ErrPostBodyEmpty = errors.New("post body is empty")
var ErrPostBodyExceedsLimit = errors.New("post body exceeds limit")
var ErrPostNotFound = errors.New("post not found")

type Service struct {
	Repository Repository
}

func (p Service) Create(ctx context.Context, post internal.Post) (internal.Post, error) {
	if post.Body == "" {
		return internal.Post{}, ErrPostBodyEmpty
	}

	if utf8.RuneCountInString(post.Body) > 140 {
		return internal.Post{}, ErrPostBodyExceedsLimit
	}

	return p.Repository.Insert(ctx, post)
}

func (s Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.Repository.Delete(ctx, id)
}

func (s Service) FindOneByID(ctx context.Context, id uuid.UUID) (internal.Post, error) {
	return s.Repository.FindOneByID(ctx, id)
}
