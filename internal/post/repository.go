package post

import (
	"context"
	"time"

	"github.com/filhodanuvem/ytgoapi/internal"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Insert(ctx context.Context, post internal.Post) (internal.Post, error)
	FindOneByID(ctx context.Context, id uuid.UUID) (internal.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type RepositoryPostgres struct {
	Conn *pgxpool.Pool
}

func (r *RepositoryPostgres) Insert(ctx context.Context, post internal.Post) (internal.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := r.Conn.QueryRow(
		ctx,
		"INSERT INTO posts (username, body) VALUES ($1, $2) RETURNING id, created_at",
		post.Username,
		post.Body).Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		return internal.Post{}, err
	}

	return post, nil
}

func (r *RepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tag, err := r.Conn.Exec(
		ctx,
		"DELETE FROM posts WHERE id = $1",
		id)

	if tag.RowsAffected() == 0 {
		return ErrPostNotFound
	}

	return err
}

func (r *RepositoryPostgres) FindOneByID(ctx context.Context, id uuid.UUID) (internal.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var post = internal.Post{ID: id}
	err := r.Conn.QueryRow(
		ctx,
		"SELECT username, body, created_at FROM posts WHERE id = $1",
		id).Scan(&post.Username, &post.Body, &post.CreatedAt)

	if err == pgx.ErrNoRows {
		return internal.Post{}, ErrPostNotFound
	}

	if err != nil {
		return internal.Post{}, err
	}

	return post, nil
}
