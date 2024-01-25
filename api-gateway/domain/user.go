package domain

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/grpc"
	"github.com/nafisalfiani/p3-ugc-7-8/api-gateway/entity"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/types/known/emptypb"
)

type user struct {
	logger     *logrus.Logger
	userClient grpc.UserServiceClient
	cache      *redis.Client
	broker     *amqp.Connection
}

type UserInterface interface {
	List(ctx context.Context) ([]entity.User, error)
	Get(ctx context.Context, filter entity.User) (entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, user entity.User) error

	UpdateUserCache(ctx context.Context, user entity.User) error
}

// initUser creates user domain
func initUser(logger *logrus.Logger, userClient grpc.UserServiceClient, cache *redis.Client, broker *amqp.Connection) UserInterface {
	return &user{
		logger:     logger,
		userClient: userClient,
		cache:      cache,
		broker:     broker,
	}
}

// List returns list of users
func (u *user) List(ctx context.Context) ([]entity.User, error) {
	userList, err := u.userClient.GetUsers(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	users := []entity.User{}
	for i := range userList.Users {
		users = append(users, entity.User{
			Id:    userList.Users[i].Id,
			Name:  userList.Users[i].Name,
			Email: userList.Users[i].Email,
		})
	}

	return users, nil
}

// Get returns specific user by email
func (u *user) Get(ctx context.Context, filter entity.User) (entity.User, error) {
	// get from cache, if no error and user found, direct return
	user, err := u.getUserCache(ctx, filter.Id)
	if err == nil && user.Id != "" {
		u.logger.Info(fmt.Sprintf("cache for user:%v found", filter.Id))
		return user, nil
	}
	u.logger.Info(fmt.Sprintf("cache for user:%v not found", filter.Id))

	res, err := u.userClient.GetUser(ctx, &grpc.User{
		Id:    filter.Id,
		Name:  filter.Name,
		Email: filter.Email,
	})
	if err != nil {
		return user, err
	}
	user.ConvertFromProto(res)

	return user, nil
}

// Create creates new data
func (u *user) Create(ctx context.Context, user entity.User) (entity.User, error) {
	res, err := u.userClient.AddUser(ctx, &grpc.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return user, err
	}

	user.ConvertFromProto(res)

	// set user cache if result found from mongo
	if err := u.setUserCache(ctx, user); err != nil {
		u.logger.Error(fmt.Sprintf("cache for user:%v failed to be set", user.Id))
	}

	return user, nil
}

// Update updates existing data
func (u *user) Update(ctx context.Context, user entity.User) (entity.User, error) {
	var newUser entity.User
	res, err := u.userClient.UpdateUser(ctx, &grpc.User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		return newUser, err
	}

	newUser.ConvertFromProto(res)

	return newUser, nil
}

// Delete deletes existing data
func (u *user) Delete(ctx context.Context, user entity.User) error {
	_, err := u.userClient.DeleteUser(ctx, &grpc.User{
		Id: user.Id,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *user) UpdateUserCache(ctx context.Context, user entity.User) error {
	return u.setUserCache(ctx, user)
}
