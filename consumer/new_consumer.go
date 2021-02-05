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

func NewConsumer(
	logger zerolog.Logger,
	stanConn stan.Conn,
	serviceName string,
	dadataClient dadata.Client,
	orgModel org.Model,
	areaModel area.Model,
	locationModel location.Model,
	managerModel manager.Model,
	okvedModel okved.Model,
	metroModel metro.Model,
) Consumer {
	return Consumer{
		logger:        logger,
		stanConn:      stanConn,
		serviceName:   serviceName,
		dadataClient:  dadataClient,
		orgModel:      orgModel,
		areaModel:     areaModel,
		locationModel: locationModel,
		managerModel:  managerModel,
		okvedModel:    okvedModel,
		metroModel:    metroModel,
		state: &state{
			done: make(chan struct{}),
		},
	}
}
