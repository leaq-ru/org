package orgimpl

import (
	"context"
	safeerr "github.com/nnqq/scr-lib-safeerr"
	"github.com/nnqq/scr-proto/codegen/go/org"
	"time"
)

func (s *server) Get(ctx context.Context, req *org.GetRequest) (res *org.GetResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	limit := uint32(20)
	if req.GetOpts().GetLimit() > 100 {
		err = safeerr.LimitOutOfRange
		return
	}
	if req.GetOpts().GetLimit() != 0 {
		limit = req.GetOpts().GetLimit()
	}

	orgs, err := s.orgModel.GetByIDs(ctx, nil, false, req.GetOpts().GetSkip(), limit)
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = safeerr.InternalServerError
		return
	}

	res, err = s.fetchShortOrgs(ctx, orgs)
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = safeerr.InternalServerError
	}
	return
}
