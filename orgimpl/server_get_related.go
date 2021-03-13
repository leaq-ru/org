package orgimpl

import (
	"context"
	safeerr "github.com/nnqq/scr-lib-safeerr"
	"github.com/nnqq/scr-org/org"
	pbOrg "github.com/nnqq/scr-proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (s *server) GetRelated(ctx context.Context, req *pbOrg.GetRelatedRequest) (res *pbOrg.GetResponse, err error) {
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

	var ids []org.ID
	areaID, e := primitive.ObjectIDFromHex(req.GetAreaId())
	if e == nil {
		ids = append(ids, org.ID{
			Val:  areaID,
			Kind: org.IDKind_areaID,
		})
	}

	okvedID, e := primitive.ObjectIDFromHex(req.GetOkvedId())
	if e == nil {
		ids = append(ids, org.ID{
			Val:  okvedID,
			Kind: org.IDKind_okvedID,
		})
	}

	excludeOrgID, e := primitive.ObjectIDFromHex(req.GetExcludeOrgId())
	if e == nil {
		ids = append(ids, org.ID{
			Val:  excludeOrgID,
			Kind: org.IDKind_excludeOrgID,
		})
	}

	orgs, err := s.orgModel.GetByIDs(ctx, ids, false, req.GetOpts().GetSkip(), limit)
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
