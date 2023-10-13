package post

import (
	"context"

	"github.com/filhodanuvem/ytgoapi/internal"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	Insert(ctx context.Context, post internal.Post) (internal.Post, error)
	FindAll(ctx context.Context) ([]*internal.Post, error)
	FindOneByID(ctx context.Context, id string) (internal.Post, error)
	Update(ctx context.Context, post internal.Post) error
	Delete(ctx context.Context, id string) error
}

type RepositoryPostgres struct {
	Conn *pgxpool.Pool
}

func (r *RepositoryPostgres) Insert(ctx context.Context, post internal.Post) (internal.Post, error) {

	err := r.Conn.QueryRow(
		ctx,
		"INSERT INTO posts (username, body) VALUES ($1, $2) RETURNING id",
		post.Username,
		post.Body).Scan(&post.ID)
	if err != nil {
		return internal.Post{}, err
	}

	return post, nil
}

func (r *RepositoryPostgres) Delete(ctx context.Context, id string) error {
	tag, err := r.Conn.Exec(
		ctx,
		"DELETE FROM posts WHERE id = $1",
		id)

	if tag.RowsAffected() == 0 {
		return ErrPostNotFound
	}

	return err
}

func (r *RepositoryPostgres) FindAll(ctx context.Context) ([]*internal.Post, error) {
	rows, err := r.Conn.Query(
		ctx,
		"SELECT id, username, body, created_at FROM posts",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []*internal.Post

	for rows.Next() {
		var item internal.Post
		if err := rows.Scan(&item.ID, &item.Username, &item.Body, &item.CreatedAt); err != nil {
			// fmt.Println(err)
			return nil, err
		}

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *RepositoryPostgres) Update(ctx context.Context, post internal.Post) error {
	_, err := r.Conn.Exec(
		ctx,
		`update "post" SET "username"= COALESCE($1, "username"), "body"= COALESCE($2, "body") where id=$3`,
		post.Username,
		post.Body,
		post.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryPostgres) FindOneByID(ctx context.Context, id string) (internal.Post, error) {

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
