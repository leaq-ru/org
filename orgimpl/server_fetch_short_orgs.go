package orgimpl

import (
	"context"
	"github.com/nnqq/scr-org/area"
	"github.com/nnqq/scr-org/location"
	"github.com/nnqq/scr-org/manager"
	"github.com/nnqq/scr-org/okved"
	"github.com/nnqq/scr-org/org"
	pbOrg "github.com/nnqq/scr-proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
)

func (s *server) fetchShortOrgs(ctx context.Context, orgs []org.Org) (res *pbOrg.GetResponse, err error) {
	unAreaIDs := map[primitive.ObjectID]struct{}{}
	unManagerIDs := map[primitive.ObjectID]struct{}{}
	unOkvedIDs := map[primitive.ObjectID]struct{}{}
	unLocationIDs := map[primitive.ObjectID]struct{}{}

	for _, o := range orgs {
		if !o.AreaID.IsZero() {
			unAreaIDs[o.AreaID] = struct{}{}
		}
		if !o.ManagerID.IsZero() {
			unManagerIDs[o.ManagerID] = struct{}{}
		}
		if !o.OkvedOsnID.IsZero() {
			unOkvedIDs[o.OkvedOsnID] = struct{}{}
		}
		if !o.LocationID.IsZero() {
			unLocationIDs[o.LocationID] = struct{}{}
		}
	}

	areaIDs := toSlice(unAreaIDs)
	managerIDs := toSlice(unManagerIDs)
	okvedIDs := toSlice(unOkvedIDs)
	locationIDs := toSlice(unLocationIDs)

	var eg errgroup.Group
	var areaDocs []area.Area
	if len(areaIDs) != 0 {
		eg.Go(func() (e error) {
			areaDocs, e = s.areaModel.GetByIDs(ctx, areaIDs)
			return
		})
	}

	var managerDocs []manager.Manager
	if len(managerIDs) != 0 {
		eg.Go(func() (e error) {
			managerDocs, e = s.managerModel.GetByIDs(ctx, managerIDs)
			return
		})
	}

	var okvedDocs []okved.Okved
	if len(okvedIDs) != 0 {
		eg.Go(func() (e error) {
			okvedDocs, e = s.okvedModel.GetByIDs(ctx, okvedIDs)
			return
		})
	}

	var locationDocs []location.Location
	if len(locationIDs) != 0 {
		eg.Go(func() (e error) {
			locationDocs, e = s.locationModel.GetByIDs(ctx, locationIDs)
			return
		})
	}
	err = eg.Wait()
	if err != nil {
		return
	}

	mArea := map[primitive.ObjectID]area.Area{}
	for _, doc := range areaDocs {
		mArea[doc.ID] = doc
	}
	mManager := map[primitive.ObjectID]manager.Manager{}
	for _, doc := range managerDocs {
		mManager[doc.ID] = doc
	}
	mOkved := map[primitive.ObjectID]okved.Okved{}
	for _, doc := range okvedDocs {
		mOkved[doc.ID] = doc
	}
	mLocation := map[primitive.ObjectID]location.Location{}
	for _, doc := range locationDocs {
		mLocation[doc.ID] = doc
	}

	res = &pbOrg.GetResponse{}
	for _, o := range orgs {
		var areaItem *pbOrg.AreaItem
		if !o.AreaID.IsZero() {

		}

		res.Orgs = append(res.Orgs, &pbOrg.OrgShort{
			Id:         o.ID.Hex(),
			Slug:       o.Slug,
			Name:       o.Name,
			Inn:        float64(o.INN),
			Kpp:        float64(o.KPP),
			Ogrn:       float64(o.OGRN),
			Kind:       pbOrg.OrgKind(o.Kind),
			Manager:    managerItem,
			Area:       areaItem,
			Location:   locationItem,
			Okved:      okvedItem,
			StatusKind: pbOrg.StatusKind(o.StatusKind),
		})
	}
	return
}

func toSlice(in map[primitive.ObjectID]struct{}) (out []primitive.ObjectID) {
	for id := range in {
		out = append(out, id)
	}
	return
}
