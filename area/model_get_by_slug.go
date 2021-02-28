package area

import (
	"context"
)

func (m Model) GetBySlug(ctx context.Context, slug string) (res Area, err error) {
	err = m.coll.FindOne(ctx, Area{
		Slug: slug,
	}).Decode(&res)
	return
}
