package post

import (
	"context"
	"sync"
	"testing"

	"github.com/filhodanuvem/ytgoapi/internal"
	"github.com/google/uuid"
)

//
// Repository Spy
//

type repositorySpy struct {
	items map[uuid.UUID]internal.Post
	mu    sync.Mutex
}

func (r *repositorySpy) Insert(ctx context.Context, post internal.Post) (internal.Post, error) {
	id := uuid.New()

	r.mu.Lock()
	defer r.mu.Unlock()
	post.ID = id
	r.items[id] = post

	return post, nil
}

func (r *repositorySpy) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.FindOneByID(ctx, id); err != nil {
		return err
	}

	delete(r.items, id)
	return nil
}

func (r *repositorySpy) FindOneByID(ctx context.Context, id uuid.UUID) (internal.Post, error) {
	post, ok := r.items[id]
	if !ok {
		return internal.Post{}, ErrPostNotFound
	}
	return post, nil
}

func (r *repositorySpy) Update(ctx context.Context, post internal.Post) error {
	postOld, err := r.FindOneByID(ctx, post.ID)
	if err != nil {
		return err
	}

	postOld.Username = post.Username
	postOld.Body = post.Body

	return nil
}

func (r *repositorySpy) FindAll(ctx context.Context) ([]internal.Post, error) {
	var items []internal.Post

	for _, v := range r.items {
		items = append(items, v)
	}

	return items, nil
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
	repo.mu = sync.Mutex{}
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
		Body:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever sinc",
		Username: "lorem",
	}
}

//
// Tests
//

func TestServiceCreate_ShouldReturnError_WhenBodyIsEmpty(t *testing.T) {
	sut := createNewService()
	post := internal.Post{}

	_, err := sut.Create(context.TODO(), post)

	if err != ErrPostBodyEmpty {
		t.Fatalf("err not assert ErrPostBodyEmpty")
	}
}

func TestServiceCreate_ShouldReturnError_WhenBodyExceedsLimit(t *testing.T) {
	sut := createNewService()
	post := internal.Post{
		Body: `Lorem I
		ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd
		psum is simply dummy t
		ext of the printing and typesetting industry. Lorem Ipsum
		has been the industry's standard dummy text ever since`,
		Username: "Lorem",
	}

	_, err := sut.Create(context.TODO(), post)
	if err != ErrPostBodyExceedsLimit {
		t.Fatalf("err not assert ErrPostBodyExceedsLimit")
	}
}

func TestServiceCreate_ShouldBeSuccessful_WhenPostPassOnValidation(t *testing.T) {
	defer repo.Clear()

	sut := createNewService()
	post := createValidPost()

	sut.Create(context.TODO(), post)

	if repo.CountEntries() != 1 {
		t.Fatalf("Invalid number of entries on repositorySpy")
	}
}

func TestServiceDelete_ShouldReturnError_WhenPostNotFound(t *testing.T) {
	defer repo.Clear()

	sut := createNewService()
	id := uuid.New()

	err := sut.Delete(context.TODO(), id)

	if err != ErrPostNotFound {
		t.Fatalf("err not assert ErrPostNotFound")
	}
}

func TestServiceDelete_ShouldBeSuccessful_WhenDeletesValidPost(t *testing.T) {
	defer repo.Clear()

	sut := createNewService()
	data := createValidPost()
	post, _ := sut.Create(context.TODO(), data)

	sut.Delete(context.TODO(), post.ID)

	if repo.CountEntries() != 0 {
		t.Fatalf("Invalid number of entries on repositorySpy")
	}
}

func TestServiceFindOneByID_ShouldReturnError_WhenPostNotFound(t *testing.T) {
	defer repo.Clear()

	sut := createNewService()
	id := uuid.New()

	_, err := sut.FindOneByID(context.TODO(), id)

	if err != ErrPostNotFound {
		t.Fatalf("err not assert ErrPostNotFound")
	}
}

func TestServiceFindOneByID_ShouldBeSuccessful_WhenDeletesValidPost(t *testing.T) {
	defer repo.Clear()

	sut := createNewService()
	data := createValidPost()
	created, _ := sut.Create(context.TODO(), data)

	post, _ := sut.FindOneByID(context.TODO(), created.ID)

	if post.ID != created.ID {
		t.Fatalf("Invalid post.ID")
	}
}
