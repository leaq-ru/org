package orgimpl

import (
	"context"
	safeerr "github.com/leaq-ru/lib-safeerr"
	"github.com/leaq-ru/org/org"
	pbOrg "github.com/leaq-ru/proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (s *server) GetByMetroId(ctx context.Context, req *pbOrg.GetByMetroIdRequest) (res *pbOrg.GetResponse, err error) {
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

	id, err := primitive.ObjectIDFromHex(req.GetMetroId())
	if err != nil {
		err = safeerr.InvalidParam("metroId")
		return
	}

	orgs, err := s.orgModel.GetByIDs(ctx, []org.ID{{
		Val:  id,
		Kind: org.IDKind_metroID,
	}}, false, req.GetOpts().GetSkip(), limit)
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
