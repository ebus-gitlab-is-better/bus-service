package biz

import (
	"context"

	"github.com/google/wire"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewBusUseCase, NewRouteUseCase, NewDriverUseCase)

type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

type RabbitData struct {
	Ch   *amqp.Channel
	Conn *amqp.Connection
}
