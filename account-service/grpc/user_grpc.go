package grpc

import (
	"context"

	"github.com/nafisalfiani/p3-ugc-7-8/account-service/entity"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/usecase"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type userGrpcServer struct {
	log  *logrus.Logger
	user usecase.UserInterface
}

func initUserGrpcServer(log *logrus.Logger, user usecase.UserInterface) *userGrpcServer {
	return &userGrpcServer{
		log:  log,
		user: user,
	}
}

func (u *userGrpcServer) mustEmbedUnimplementedUserServiceServer() {}

func (u *userGrpcServer) GetUser(ctx context.Context, req *User) (*User, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		u.log.Error(err)
	}

	user, err := u.user.Get(ctx, entity.User{
		Id:    id,
		Name:  req.GetName(),
		Email: req.GetEmail(),
	})
	if err != nil {
		return nil, err
	}

	res := &User{
		Id:       user.Id.Hex(),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	return res, nil
}

func (u *userGrpcServer) AddUser(ctx context.Context, req *User) (*User, error) {
	newUser, err := u.user.Create(ctx, entity.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	res := &User{
		Id:    newUser.Id.Hex(),
		Name:  newUser.Name,
		Email: newUser.Email,
	}

	return res, nil
}

func (u *userGrpcServer) UpdateUser(ctx context.Context, req *User) (*User, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	user, err := u.user.Update(ctx, entity.User{
		Id:    id,
		Name:  req.GetName(),
		Email: req.GetEmail(),
	})
	if err != nil {
		return nil, err
	}
	u.log.Debug(user)

	res := &User{
		Id:    user.Id.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}

	return res, nil
}

func (u *userGrpcServer) DeleteUser(ctx context.Context, req *User) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	if err := u.user.Delete(ctx, entity.User{
		Id: id,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (u *userGrpcServer) GetUsers(ctx context.Context, in *emptypb.Empty) (*UserList, error) {
	users, err := u.user.List(ctx)
	if err != nil {
		return nil, err
	}

	res := &UserList{}
	for i := range users {
		res.Users = append(res.Users, &User{
			Id:    users[i].Id.Hex(),
			Name:  users[i].Name,
			Email: users[i].Email,
		})
	}

	return res, nil
}
