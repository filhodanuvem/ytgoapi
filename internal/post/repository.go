package post

import (
	"context"
	"fmt"
	"time"

	"github.com/filhodanuvem/ytgoapi/internal"
	"github.com/google/uuid"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	Conn *pgxpool.Pool
}

func (r *Repository) Insert(post internal.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println(r.Conn)

	_, err := r.Conn.Exec(
		ctx,
		"INSERT INTO posts (username, body) VALUES ($1, $2)",
		post.Username,
		post.Body)

	return err
}

func (r *Repository) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func (r *Repository) FindOneByID(id uuid.UUID) (internal.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var post internal.Post
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
