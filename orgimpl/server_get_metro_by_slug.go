package orgimpl

import (
	"context"
	"errors"
	safeerr "github.com/leaq-ru/lib-safeerr"
	pbOrg "github.com/leaq-ru/proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (s *server) GetMetroBySlug(
	ctx context.Context,
	req *pbOrg.GetMetroBySlugRequest,
) (
	res *pbOrg.GetMetroBySlugResponse,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	met, err := s.metroModel.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = safeerr.NotFound("metro")
		} else {
			s.logger.Error().Err(err).Send()
			err = safeerr.InternalServerError
		}
		return
	}

	var area *pbOrg.AreaFullItem
	if !met.AreaID.IsZero() {
		ar, e := s.areaModel.GetByID(ctx, met.AreaID)
		if e != nil {
			s.logger.Error().Err(e).Send()
			err = safeerr.InternalServerError
			return
		}

		area = &pbOrg.AreaFullItem{
			Id:       ar.ID.Hex(),
			Slug:     ar.Slug,
			Name:     ar.Name,
			FiasId:   ar.FiasID,
			KladrId:  ar.KladrID,
			Type:     ar.Type,
			TypeFull: ar.TypeFull,
		}
	}

	res = &pbOrg.GetMetroBySlugResponse{
		Metro: &pbOrg.MetroFullItem{
			Id:   met.ID.Hex(),
			Slug: met.Slug,
			Name: met.Name,
			Line: met.Line,
			Area: area,
		},
	}
	return
}
