package usecase

import (
	"context"

	"github.com/nafisalfiani/p3-ugc-7-8/api-gateway/config"
	"github.com/nafisalfiani/p3-ugc-7-8/api-gateway/domain"
	"github.com/nafisalfiani/p3-ugc-7-8/api-gateway/entity"
)

type user struct {
	cfg  *config.Value
	user domain.UserInterface
}

type UserInterface interface {
	List(ctx context.Context) ([]entity.User, error)
	Get(ctx context.Context, filter entity.User) (entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, user entity.User) error

	UpdateUserCache(ctx context.Context, user entity.User) error
}

// initUser creates user repository
func initUser(cfg *config.Value, userDom domain.UserInterface) UserInterface {
	return &user{
		cfg:  cfg,
		user: userDom,
	}
}

func (u *user) List(ctx context.Context) ([]entity.User, error) {
	return u.user.List(ctx)
}

func (u *user) Get(ctx context.Context, filter entity.User) (entity.User, error) {
	return u.user.Get(ctx, filter)
}

func (u *user) Create(ctx context.Context, user entity.User) (entity.User, error) {
	return u.user.Create(ctx, user)
}

func (u *user) Update(ctx context.Context, user entity.User) (entity.User, error) {
	return u.user.Update(ctx, user)
}
func (u *user) Delete(ctx context.Context, user entity.User) error {
	return u.user.Delete(ctx, user)
}

func (u *user) UpdateUserCache(ctx context.Context, user entity.User) error {
	return u.user.UpdateUserCache(ctx, user)
}
