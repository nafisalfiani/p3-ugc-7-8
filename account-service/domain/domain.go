package domain

import (
	"github.com/go-redis/redis/v8"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Domains struct {
	User UserInterface
}

func Init(db *mongo.Client, cache *redis.Client, logger *logrus.Logger) *Domains {
	return &Domains{
		User: initUser(logger, db.Database("account-service").Collection("user"), cache),
	}
}

func errorAlias(err error) error {
	// Check if the error is a mongo duplicate key error
	mongoErr, ok := err.(mongo.WriteException)
	if ok {
		for _, writeError := range mongoErr.WriteErrors {
			if writeError.Code == 11000 {
				return errors.ErrDuplicatedKey
			}
		}
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound), errors.Is(err, mongo.ErrNoDocuments):
		return errors.ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return errors.ErrDuplicatedKey
	default:
		return err
	}
}
