package consumer

import (
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-org/dadata"
	"github.com/rs/zerolog"
)

type state struct {
	sub               stan.Subscription
	subscribeCalledOK bool
	drain             bool
	done              chan struct{}
}

type Consumer struct {
	logger      zerolog.Logger
	stanConn    stan.Conn
	serviceName string
	tokenPool   dadata.TokenPool
	state       *state
}
