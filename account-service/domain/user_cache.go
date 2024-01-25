package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nafisalfiani/p3-ugc-7-8/account-service/entity"
)

func (s *user) getUserCache(ctx context.Context, userId string) (entity.User, error) {
	var user entity.User
	userStr, err := s.cache.Get(ctx, fmt.Sprintf(entity.RedisKeyUser, userId)).Result()
	if err != nil {
		return user, err
	}

	if err := json.Unmarshal([]byte(userStr), &user); err != nil {
		return user, err
	}

	return user, nil
}

func (s *user) setUserCache(ctx context.Context, user entity.User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := s.cache.Set(ctx, fmt.Sprintf(entity.RedisKeyUser, user.Id.Hex()), userJson, time.Hour).Err(); err != nil {
		return err
	}

	return nil
}
