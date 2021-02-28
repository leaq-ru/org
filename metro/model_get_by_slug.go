package metro

import (
	"context"
)

func (m Model) GetBySlug(ctx context.Context, slug string) (res Metro, err error) {
	err = m.coll.FindOne(ctx, Metro{
		Slug: slug,
	}).Decode(&res)
	return
}
