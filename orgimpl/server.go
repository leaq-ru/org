package orgimpl

import (
	"github.com/nnqq/scr-org/area"
	"github.com/nnqq/scr-org/location"
	"github.com/nnqq/scr-org/manager"
	"github.com/nnqq/scr-org/metro"
	"github.com/nnqq/scr-org/okved"
	"github.com/nnqq/scr-org/org"
	pbOrg "github.com/nnqq/scr-proto/codegen/go/org"
	"github.com/rs/zerolog"
)

type server struct {
	pbOrg.UnimplementedOrgServer
	logger        zerolog.Logger
	orgModel      org.Model
	areaModel     area.Model
	locationModel location.Model
	managerModel  manager.Model
	metroModel    metro.Model
	okvedModel    okved.Model
}
