package orgimpl

import (
	"github.com/leaq-ru/org/area"
	"github.com/leaq-ru/org/location"
	"github.com/leaq-ru/org/manager"
	"github.com/leaq-ru/org/metro"
	"github.com/leaq-ru/org/okved"
	"github.com/leaq-ru/org/org"
	"github.com/rs/zerolog"
)

func NewServer(
	logger zerolog.Logger,
	orgModel org.Model,
	areaModel area.Model,
	locationModel location.Model,
	managerModel manager.Model,
	okvedModel okved.Model,
	metroModel metro.Model,
) *server {
	return &server{
		logger:        logger,
		orgModel:      orgModel,
		areaModel:     areaModel,
		locationModel: locationModel,
		managerModel:  managerModel,
		metroModel:    metroModel,
		okvedModel:    okvedModel,
	}
}
