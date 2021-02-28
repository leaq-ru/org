package orgimpl

import (
	"context"
	"errors"
	safeerr "github.com/nnqq/scr-lib-safeerr"
	pbOrg "github.com/nnqq/scr-proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (s *server) GetManagerBySlug(
	ctx context.Context,
	req *pbOrg.GetManagerBySlugRequest,
) (
	res *pbOrg.GetManagerBySlugResponse,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	man, err := s.managerModel.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = safeerr.NotFound("manager")
		} else {
			s.logger.Error().Err(err).Send()
			err = safeerr.InternalServerError
		}
		return
	}

	res = &pbOrg.GetManagerBySlugResponse{
		Manager: &pbOrg.ManagerItem{
			Id:   man.ID.Hex(),
			Slug: man.Slug,
			Name: man.Name,
		},
	}
	return
}
