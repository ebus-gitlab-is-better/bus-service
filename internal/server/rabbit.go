package server

import (
	"bus-service/internal/biz"
	"bus-service/pkg/rabbit"
	"context"
	"encoding/json"
	"log"
)

func NewRabbitConn(ch *biz.RabbitData, uc *biz.RouteUseCase) *rabbit.RabbitConn {

	msgs, err := ch.Ch.Consume(
		"accident",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Не удалось зарегистрировать consumer: %s", err)
	}

	go func() {
		for d := range msgs {
			var accident biz.Accident
			err := json.Unmarshal(d.Body, &accident)
			if err == nil {
				uc.NewAccident(context.TODO(), &accident)
			}
		}
	}()

	return rabbit.NewRabbitConn(ch.Conn, ch.Ch)
}
