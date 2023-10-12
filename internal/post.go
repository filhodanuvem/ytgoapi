package internal

import (
	"errors"
	"time"
)

type Post struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type ParamsUpdatePost struct {
	ID       string
	Username string
	Body     string
}

func (p *ParamsUpdatePost) Validate() (Post, error) {
	if p.ID == "" {
		return Post{}, errors.New("id empty")
	}

	if utf8.RuneCountInString(p.Body) > 140 {
		return internal.Post{}, ErrPostBodyExceedsLimit
	}

	return Post{
		ID:       p.ID,
		Username: p.Username,
		Body:     p.Body,
	}, nil
}
