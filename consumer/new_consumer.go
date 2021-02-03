package consumer

import (
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-org/dadata"
	"github.com/rs/zerolog"
)

func NewConsumer(
	logger zerolog.Logger,
	stanConn stan.Conn,
	serviceName string,
	tokenPool dadata.TokenPool,
) Consumer {
	return Consumer{
		logger:      logger,
		stanConn:    stanConn,
		serviceName: serviceName,
		tokenPool:   tokenPool,
		state: &state{
			done: make(chan struct{}),
		},
	}
}
