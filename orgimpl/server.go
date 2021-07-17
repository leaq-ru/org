package orgimpl

import (
	"github.com/leaq-ru/org/area"
	"github.com/leaq-ru/org/location"
	"github.com/leaq-ru/org/manager"
	"github.com/leaq-ru/org/metro"
	"github.com/leaq-ru/org/okved"
	"github.com/leaq-ru/org/org"
	pbOrg "github.com/leaq-ru/proto/codegen/go/org"
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
