package location

import (
	"context"
)

func (m Model) GetBySlug(ctx context.Context, slug string) (res Location, err error) {
	err = m.coll.FindOne(ctx, Location{
		Slug: slug,
	}).Decode(&res)
	return
}
