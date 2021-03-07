package orgimpl

import (
	"context"
	"errors"
	safeerr "github.com/nnqq/scr-lib-safeerr"
	pbOrg "github.com/nnqq/scr-proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"strings"
	"time"
)

func (s *server) GetBySlug(ctx context.Context, req *pbOrg.GetBySlugRequest) (res *pbOrg.GetBySlugResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	slugSlice := strings.Split(req.GetSlug(), "-")
	lastIndex := 0
	if len(slugSlice) != 0 {
		lastIndex = len(slugSlice) - 1
	}

	inn, err := strconv.Atoi(slugSlice[lastIndex])
	if err != nil {
		err = safeerr.NotFound("org")
		return
	}

	orgs, err := s.orgModel.GetByINN(ctx, uint64(inn))
	if err != nil {
		s.logger.Error().Err(err).Send()
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = safeerr.NotFound("org")
		} else {
			err = safeerr.InternalServerError
		}
		return
	}

	res, err = s.fetchOrgWithBranchesAndRelated(ctx, orgs)
	if err != nil {
		s.logger.Error().Err(err).Send()
		err = safeerr.InternalServerError
	}
	return
}
