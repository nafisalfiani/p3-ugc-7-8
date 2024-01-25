package domain

import (
	"encoding/json"
	"fmt"

	"github.com/nafisalfiani/p3-ugc-7-8/account-service/entity"
	"github.com/streadway/amqp"
)

func (u *user) publishUser(user entity.User) error {
	ch, err := u.broker.Channel()
	if err != nil {
		u.logger.Error(fmt.Sprintf("failed to publish message for user:%v", user.Id))
		return nil
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(entity.QueueTopicUserAdded, false, false, false, false, nil)
	if err != nil {
		return err
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        userJson,
	}

	return ch.Publish("", queue.Name, false, false, message)
}
