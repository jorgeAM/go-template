package infrastructure

import (
	"context"
	"errors"
	"sync"

	"github.com/jorgeAM/base-api/internal/user"
)

var _ user.UserRepository = (*InMemUserRepository)(nil)

type InMemUserRepository struct {
	mux sync.RWMutex

	items map[string]*user.User
}

func NewInMemUserRepository() *InMemUserRepository {
	return &InMemUserRepository{}
}

func (i *InMemUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	user, ok := i.items[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (i *InMemUserRepository) Save(ctx context.Context, user *user.User) error {
	i.mux.Lock()
	defer i.mux.Unlock()

	i.items[user.ID] = user

	return nil
}
