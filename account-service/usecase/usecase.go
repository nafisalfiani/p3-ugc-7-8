package usecase

import (
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/config"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/domain"

	"github.com/sirupsen/logrus"
)

type Usecases struct {
	User UserInterface
}

func Init(cfg *config.Value, logger *logrus.Logger, dom *domain.Domains) *Usecases {
	return &Usecases{
		User: initUser(cfg, dom.User),
	}
}
