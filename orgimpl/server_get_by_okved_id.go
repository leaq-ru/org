package orgimpl

import (
	"context"
	safeerr "github.com/nnqq/scr-lib-safeerr"
	"github.com/nnqq/scr-org/org"
	pbOrg "github.com/nnqq/scr-proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (s *server) GetByOkvedId(ctx context.Context, req *pbOrg.GetByOkvedIdRequest) (res *pbOrg.GetResponse, err error) {
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

	id, err := primitive.ObjectIDFromHex(req.GetOkvedId())
	if err != nil {
		err = safeerr.InvalidParam("okvedId")
		return
	}

	orgs, err := s.orgModel.GetByIDs(ctx, []org.ID{{
		Val:  id,
		Kind: org.IDKind_okvedID,
	}}, req.GetOpts().GetSkip(), limit)
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
