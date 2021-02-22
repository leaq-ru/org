package orgimpl

import (
	"github.com/nnqq/scr-org/area"
	"github.com/nnqq/scr-org/location"
	"github.com/nnqq/scr-org/manager"
	"github.com/nnqq/scr-org/metro"
	"github.com/nnqq/scr-org/okved"
	"github.com/nnqq/scr-org/org"
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
