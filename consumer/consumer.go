package consumer

import (
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-org/area"
	"github.com/nnqq/scr-org/dadata"
	"github.com/nnqq/scr-org/location"
	"github.com/nnqq/scr-org/manager"
	"github.com/nnqq/scr-org/metro"
	"github.com/nnqq/scr-org/okved"
	"github.com/nnqq/scr-org/org"
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
