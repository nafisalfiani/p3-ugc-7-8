package domain

import (
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/grpc"

	"github.com/sirupsen/logrus"
)

type Domains struct {
	User UserInterface
}

func Init(logger *logrus.Logger, userClient grpc.UserServiceClient) *Domains {
	return &Domains{
		User: initUser(logger, userClient),
	}
}
