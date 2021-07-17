package consumer

import (
	"github.com/leaq-ru/org/area"
	"github.com/leaq-ru/org/dadata"
	"github.com/leaq-ru/org/location"
	"github.com/leaq-ru/org/manager"
	"github.com/leaq-ru/org/metro"
	"github.com/leaq-ru/org/okved"
	"github.com/leaq-ru/org/org"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
)

type state struct {
	sub               stan.Subscription
	subscribeCalledOK bool
	drain             bool
	done              chan struct{}
}

type Consumer struct {
	logger        zerolog.Logger
	stanConn      stan.Conn
	serviceName   string
	dadataClient  dadata.Client
	orgModel      org.Model
	areaModel     area.Model
	locationModel location.Model
	managerModel  manager.Model
	okvedModel    okved.Model
	metroModel    metro.Model
	state         *state
}
