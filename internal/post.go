package internal

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type ParamsUpdatePost struct {
	ID       string
	Username string
	Body     string
}

var (
	ErrPostBodyExceedsLimit = errors.New("post body exceeds limit")
	ErrIdEmpty              = errors.New("id empty")
	ErrUUIDInvalid          = errors.New("uuid invalid")
)

func (p *ParamsUpdatePost) Validate() (Post, error) {
	if p.ID == "" {
		return Post{}, ErrUUIDInvalid
	}

	idParsed, err := uuid.Parse(p.ID)
	if err != nil {
		return Post{}, ErrUUIDInvalid
	}

	if utf8.RuneCountInString(p.Body) > 140 {
		return Post{}, ErrPostBodyExceedsLimit
	}

	return Post{
		ID:       idParsed,
		Username: p.Username,
		Body:     p.Body,
	}, nil
}
