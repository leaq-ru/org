package consumer

import (
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
)

func NewConsumer(
	logger zerolog.Logger,
	stanConn stan.Conn,
	serviceName string,
) Consumer {
	return Consumer{
		logger:      logger,
		stanConn:    stanConn,
		serviceName: serviceName,
		state: &state{
			done: make(chan struct{}),
		},
	}
}
