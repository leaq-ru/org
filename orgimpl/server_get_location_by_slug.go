package orgimpl

import (
	"context"
	"errors"
	safeerr "github.com/nnqq/scr-lib-safeerr"
	pbOrg "github.com/nnqq/scr-proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (s *server) GetLocationBySlug(
	ctx context.Context,
	req *pbOrg.GetLocationBySlugRequest,
) (
	res *pbOrg.GetLocationBySlugResponse,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	loc, err := s.locationModel.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = safeerr.NotFound("location")
		} else {
			s.logger.Error().Err(err).Send()
			err = safeerr.InternalServerError
		}
		return
	}

	res = &pbOrg.GetLocationBySlugResponse{
		Location: &pbOrg.LocationItem{
			Id:   loc.ID.Hex(),
			Slug: loc.Slug,
			Name: loc.Name,
		},
	}
	return
}
