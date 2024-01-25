package handler

import (
	"context"
	"encoding/json"

	"github.com/nafisalfiani/p3-ugc-7-8/api-gateway/entity"
	"github.com/streadway/amqp"
)

func (h *Handler) StartConsumer() (*amqp.Channel, error) {
	ch, err := h.broker.Channel()
	if err != nil {
		return ch, err
	}

	queue, err := ch.QueueDeclare(entity.QueueTopicUserAdded, false, false, false, false, nil)
	if err != nil {
		return ch, err
	}

	// RabbitMQ consumer
	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return ch, err
	}

	for d := range msgs {
		h.logger.Info("received new message from consumer")

		var user entity.User
		if err := json.Unmarshal(d.Body, &user); err != nil {
			h.logger.Error(err)
			continue
		}

		if err := h.user.UpdateUserCache(context.Background(), user); err != nil {
			h.logger.Error(err)
			continue
		}
	}

	return ch, nil
}
