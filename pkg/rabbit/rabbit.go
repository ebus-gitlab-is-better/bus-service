package rabbit

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConn struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitConn(conn *amqp.Connection, channel *amqp.Channel) *RabbitConn {
	return &RabbitConn{conn: conn, channel: channel}
}

func (s *RabbitConn) Start(ctx context.Context) error {
	return nil
}

func (s *RabbitConn) Stop(ctx context.Context) error {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
	return nil
}
