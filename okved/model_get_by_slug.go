package okved

import (
	"context"
)

func (m Model) GetBySlug(ctx context.Context, slug string) (res Okved, err error) {
	err = m.coll.FindOne(ctx, Okved{
		Slug: slug,
	}).Decode(&res)
	return
}
