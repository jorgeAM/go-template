package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jorgeAM/base-api/internal/user/domain"
)

var _ domain.UserRepository = (*InMemUserRepository)(nil)

type InMemUserRepository struct {
	mux sync.RWMutex

	items map[string]*domain.User
}

func NewInMemUserRepository() *InMemUserRepository {
	return &InMemUserRepository{
		mux:   sync.RWMutex{},
		items: map[string]*domain.User{},
	}
}

func (i *InMemUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user, ok := i.items[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (i *InMemUserRepository) Save(ctx context.Context, user *domain.User) error {
	i.mux.Lock()
	defer i.mux.Unlock()

	i.items[user.ID] = user

	fmt.Println(i.items)

	return nil
}
