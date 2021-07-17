package orgimpl

import (
	"context"
	"errors"
	safeerr "github.com/leaq-ru/lib-safeerr"
	pbOrg "github.com/leaq-ru/proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (s *server) GetOkvedBySlug(
	ctx context.Context,
	req *pbOrg.GetOkvedBySlugRequest,
) (
	res *pbOrg.GetOkvedBySlugResponse,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	okv, err := s.okvedModel.GetBySlug(ctx, req.GetSlug())
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = safeerr.NotFound("okved")
		} else {
			s.logger.Error().Err(err).Send()
			err = safeerr.InternalServerError
		}
		return
	}

	res = &pbOrg.GetOkvedBySlugResponse{
		Okved: &pbOrg.OkvedFullItem{
			Id:           okv.ID.Hex(),
			Slug:         okv.Slug,
			Name:         okv.Name,
			Code:         okv.Code,
			CodeWithName: okv.CodeWithName,
			Kind:         pbOrg.OkvedKind(okv.Kind),
		},
	}
	return
}
