package manager

import (
	"context"
)

func (m Model) GetBySlug(ctx context.Context, slug string) (res Manager, err error) {
	err = m.coll.FindOne(ctx, Manager{
		Slug: slug,
	}).Decode(&res)
	return
}
