package entity

import "github.com/nafisalfiani/p3-ugc-7-8/account-service/grpc"

const (
	QueueTopicUser      = "user"
	QueueTopicUserAdded = "user_added"

	RedisKeyUser = "api-gw:user:%v"
)

type User struct {
	Id       string `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name,omitempty"`
	Email    string `json:"email" bson:"email,omitempty"`
	Password string `json:"-" bson:"password,omitempty"`
}

func (u *User) ConvertFromProto(user *grpc.User) {
	u.Id = user.GetId()
	u.Name = user.GetName()
	u.Email = user.GetEmail()
	u.Password = user.GetPassword()
}

type UserCreateRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type UserUpdateRequest struct {
	Id    string `param:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserGetRequest struct {
	Id string `param:"id" validate:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResp struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
