package domain

import (
	"api-gateway/entity"
	"context"

	"github.com/nafisalfiani/p3-ugc-7-8/account-service/grpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type user struct {
	logger     *logrus.Logger
	userClient grpc.UserServiceClient
}

type UserInterface interface {
	List(ctx context.Context) ([]entity.User, error)
	Get(ctx context.Context, filter entity.User) (entity.User, error)
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, user entity.User) error
}

// initUser creates user domain
func initUser(logger *logrus.Logger, userClient grpc.UserServiceClient) UserInterface {
	return &user{
		logger:     logger,
		userClient: userClient,
	}
}

// List returns list of users
func (s *user) List(ctx context.Context) ([]entity.User, error) {
	userList, err := s.userClient.GetUsers(ctx, &emptypb.Empty{})
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
func (s *user) Get(ctx context.Context, filter entity.User) (entity.User, error) {
	var user entity.User
	res, err := s.userClient.GetUser(ctx, &grpc.User{
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
func (s *user) Create(ctx context.Context, user entity.User) (entity.User, error) {
	res, err := s.userClient.AddUser(ctx, &grpc.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return user, err
	}

	user.ConvertFromProto(res)

	return user, nil
}

// Update updates existing data
func (s *user) Update(ctx context.Context, user entity.User) (entity.User, error) {
	var newUser entity.User
	res, err := s.userClient.UpdateUser(ctx, &grpc.User{
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
func (s *user) Delete(ctx context.Context, user entity.User) error {
	_, err := s.userClient.DeleteUser(ctx, &grpc.User{
		Id: user.Id,
	})
	if err != nil {
		return err
	}

	return nil
}
