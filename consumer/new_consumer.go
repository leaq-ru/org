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
