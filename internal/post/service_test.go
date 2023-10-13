package post

import (
	"context"
	"testing"

	"github.com/filhodanuvem/ytgoapi/internal"
	"github.com/google/uuid"
)

//
// Repository Spy
//

type repositorySpy struct {
	items map[string]internal.Post
}

func (r *repositorySpy) Insert(post internal.Post) (internal.Post, error) {
	id := uuid.NewString()

	post.ID = id
	r.items[id] = post

	return post, nil
}

func (r *repositorySpy) Delete(id uuid.UUID) error {
	if _, err := r.FindOneByID(id); err != nil {
		return err
	}

	delete(r.items, id)
	return nil
}

func (r *repositorySpy) FindOneByID(id uuid.UUID) (internal.Post, error) {
	post, ok := r.items[id]
	if !ok {
		return internal.Post{}, ErrPostNotFound
	}
	return post, nil
}

func (r *repositorySpy) CountEntries() int {
	return len(r.items)
}

func (r *repositorySpy) Clear() {
	clear(r.items)
}

//
// Setup
//

func createRepository() *repositorySpy {
	repo := repositorySpy{}
	repo.items = make(map[string]internal.Post)
	return &repo
}

var repo = createRepository()

//
// Utils
//

func createNewService() Service {
	return Service{
		Repository: repo,
	}
}

func createValidPost() internal.Post {
	return internal.Post{
		Body: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever sinc",
	}
}

//
// Tests
//

func TestServiceCreate_ShouldReturnError_WhenBodyIsEmpty(t *testing.T) {
	sut := createNewService()
	post := internal.Post{}

	ctx := context.Background()

	_, err := sut.Create(ctx, post)

	if err != ErrPostBodyEmpty {
		t.Fatalf("err not assert ErrPostBodyEmpty")
	}
}

func TestServiceCreate_ShouldReturnError_WhenBodyExceedsLimit(t *testing.T) {
	sut := createNewService()
	post := internal.Post{
		Body: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since",
	}

	ctx := context.Background()

	_, err := sut.Create(ctx, post)

	if err != ErrPostBodyExceedsLimit {
		t.Fatalf("err not assert ErrPostBodyExceedsLimit")
	}
}

func TestServiceCreate_ShouldBeSuccessful_WhenPostPassOnValidation(t *testing.T) {
	defer repo.Clear()

	sut := createNewService()
	post := createValidPost()

	ctx := context.Background()

	sut.Create(ctx, post)

	if repo.CountEntries() != 1 {
		t.Fatalf("Invalid number of entries on repositorySpy")
	}
}

func TestServiceDelete_ShouldReturnError_WhenPostNotFound(t *testing.T) {
	defer repo.Clear()

	sut := createNewService()
	id := uuid.NewString()
	ctx := context.Background()

	err := sut.Delete(ctx, id)

	if err != ErrPostNotFound {
		t.Fatalf("err not assert ErrPostNotFound")
	}
}

func TestServiceDelete_ShouldBeSuccessful_WhenDeletesValidPost(t *testing.T) {
	defer repo.Clear()

	sut := createNewService()
	data := createValidPost()
	post, _ := sut.Create(data)

	ctx := context.Background()

	sut.Delete(ctx, post.ID)

	if repo.CountEntries() != 0 {
		t.Fatalf("Invalid number of entries on repositorySpy")
	}
}

func TestServiceFindOneByID_ShouldReturnError_WhenPostNotFound(t *testing.T) {
	defer repo.Clear()

	sut := createNewService()
	id := uuid.NewString()

	ctx := context.Background()

	_, err := sut.FindOneByID(ctx, id)

	if err != ErrPostNotFound {
		t.Fatalf("err not assert ErrPostNotFound")
	}
}

func TestServiceFindOneByID_ShouldBeSuccessful_WhenDeletesValidPost(t *testing.T) {
	defer repo.Clear()

	sut := createNewService()
	data := createValidPost()
	created, _ := sut.Create(data)

	ctx := context.Background()

	post, _ := sut.FindOneByID(ctx, created.ID)

	if post.ID != created.ID {
		t.Fatalf("Invalid post.ID")
	}
}
