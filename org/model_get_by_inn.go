package org

import "context"

func (m Model) GetByINN(ctx context.Context, inn uint64) (res []Org, err error) {
	cur, err := m.coll.Find(ctx, Org{
		INN: inn,
	})
	if err != nil {
		return
	}

	err = cur.All(ctx, &res)
	return
}
