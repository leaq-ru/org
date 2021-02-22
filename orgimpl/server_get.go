package orgimpl

import (
	"context"
	"errors"
	"github.com/nnqq/scr-proto/codegen/go/org"
	"net/http"
	"time"
)

func (s *server) Get(ctx context.Context, req *org.GetRequest) (res *org.GetResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ise := errors.New(http.StatusText(http.StatusInternalServerError))

	limit := uint32(20)
	if req.GetOpts().GetLimit() > 100 {
		err = errors.New("limit too large")
		return
	}
	if req.GetOpts().GetLimit() != 0 {
		limit = req.GetOpts().GetLimit()
	}

	orgs, err := s.orgModel.GetByIDs(ctx, nil, req.GetOpts().GetSkip(), limit)
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = ise
		return
	}

	res, err = s.fetchShortOrgs(ctx, orgs)
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = ise
	}
	return
}
