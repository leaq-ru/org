package orgimpl

import (
	"context"
	"errors"
	safeerr "github.com/leaq-ru/lib-safeerr"
	pbOrg "github.com/leaq-ru/proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (s *server) GetAreaBySlug(
	ctx context.Context,
	req *pbOrg.GetAreaBySlugRequest,
) (
	res *pbOrg.GetAreaBySlugResponse,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ar, err := s.areaModel.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = safeerr.NotFound("area")
		} else {
			s.logger.Error().Err(err).Send()
			err = safeerr.InternalServerError
		}
		return
	}

	res = &pbOrg.GetAreaBySlugResponse{
		Area: &pbOrg.AreaFullItem{
			Id:       ar.ID.Hex(),
			Slug:     ar.Slug,
			Name:     ar.Name,
			FiasId:   ar.FiasID,
			KladrId:  ar.KladrID,
			Type:     ar.Type,
			TypeFull: ar.TypeFull,
		},
	}
	return
}
