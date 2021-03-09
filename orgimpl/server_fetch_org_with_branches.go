package orgimpl

import (
	"context"
	"errors"
	"github.com/nnqq/scr-org/area"
	"github.com/nnqq/scr-org/location"
	"github.com/nnqq/scr-org/manager"
	"github.com/nnqq/scr-org/metro"
	"github.com/nnqq/scr-org/okved"
	"github.com/nnqq/scr-org/org"
	"github.com/nnqq/scr-proto/codegen/go/opts"
	pbOrg "github.com/nnqq/scr-proto/codegen/go/org"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

func (s *server) fetchOrgWithBranchesAndRelated(
	ctx context.Context,
	orgs []org.Org,
) (
	res *pbOrg.GetBySlugResponse,
	err error,
) {
	unAreaIDs := map[primitive.ObjectID]struct{}{}
	unManagerIDs := map[primitive.ObjectID]struct{}{}
	unOkvedIDs := map[primitive.ObjectID]struct{}{}
	unLocationIDs := map[primitive.ObjectID]struct{}{}
	unMetroIDs := map[primitive.ObjectID]struct{}{}

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
		for _, m := range o.Metros {
			if !m.ID.IsZero() {
				unMetroIDs[m.ID] = struct{}{}
			}
		}
		for _, id := range o.OkvedDopIDs {
			unOkvedIDs[id] = struct{}{}
		}
	}

	areaIDs := toSlice(unAreaIDs)
	managerIDs := toSlice(unManagerIDs)
	okvedIDs := toSlice(unOkvedIDs)
	locationIDs := toSlice(unLocationIDs)
	metroIDs := toSlice(unMetroIDs)

	var eg errgroup.Group
	var (
		areaMu   sync.Mutex
		areaDocs []area.Area
	)
	if len(areaIDs) != 0 {
		eg.Go(func() error {
			ad, e := s.areaModel.GetByIDs(ctx, areaIDs)
			if e != nil {
				return e
			}

			areaMu.Lock()
			areaDocs = append(areaDocs, ad...)
			areaMu.Unlock()
			return nil
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

	var metroDocs []metro.Metro
	if len(metroIDs) != 0 {
		eg.Go(func() error {
			md, e := s.metroModel.GetByIDs(ctx, metroIDs)
			if e != nil {
				return e
			}
			metroDocs = md

			var dopAreaIDs []primitive.ObjectID
			for _, doc := range metroDocs {
				dopAreaIDs = append(dopAreaIDs, doc.AreaID)
			}

			ad, e := s.areaModel.GetByIDs(ctx, dopAreaIDs)
			if e != nil {
				return e
			}

			areaMu.Lock()
			areaDocs = append(areaDocs, ad...)
			areaMu.Unlock()
			return nil
		})
	}

	var related []*pbOrg.OrgShort
	eg.Go(func() error {
		var mainOrg org.Org
		for _, o := range orgs {
			if o.BranchKind == org.BranchKind_main {
				mainOrg = o
				break
			}
		}
		if mainOrg.ID.IsZero() {
			return nil
		}

		rel, e := s.GetRelated(ctx, &pbOrg.GetRelatedRequest{
			Opts: &opts.SkipLimit{
				Limit: 6,
			},
			AreaId:       mainOrg.AreaID.Hex(),
			OkvedId:      mainOrg.OkvedOsnID.Hex(),
			ExcludeOrgId: mainOrg.ID.Hex(),
		})
		if e != nil {
			return e
		}

		related = rel.GetOrgs()
		return nil
	})
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
	mMetro := map[primitive.ObjectID]metro.Metro{}
	for _, doc := range metroDocs {
		mMetro[doc.ID] = doc
	}

	res = &pbOrg.GetBySlugResponse{}
	for _, o := range orgs {
		var (
			areaItem     *pbOrg.AreaItem
			areaFullItem *pbOrg.AreaFullItem
		)
		if !o.AreaID.IsZero() {
			val, ok := mArea[o.AreaID]
			if !ok {
				err = errors.New("expected to get area from map, but nothing found o.AreaID=" + o.AreaID.Hex())
				return
			}

			areaItem = &pbOrg.AreaItem{
				Id:   val.ID.Hex(),
				Slug: val.Slug,
				Name: val.Name,
			}

			areaFullItem = &pbOrg.AreaFullItem{
				Id:       val.ID.Hex(),
				Slug:     val.Slug,
				Name:     val.Name,
				FiasId:   val.FiasID,
				KladrId:  val.KladrID,
				Type:     val.Type,
				TypeFull: val.TypeFull,
			}
		}

		var locationItem *pbOrg.LocationItem
		if !o.AreaID.IsZero() {
			val, ok := mLocation[o.LocationID]
			if !ok {
				err = errors.New("expected to get location from map, but nothing found o.LocationID=" + o.LocationID.Hex())
				return
			}

			locationItem = &pbOrg.LocationItem{
				Id:   val.ID.Hex(),
				Slug: val.Slug,
				Name: val.Name,
			}
		}

		if o.BranchKind == org.BranchKind_branch {
			res.Branches = append(res.Branches, &pbOrg.Branch{
				Name:       o.Name,
				Area:       areaItem,
				Location:   locationItem,
				StatusKind: pbOrg.StatusKind(o.StatusKind),
			})
			continue
		}

		var managerItem *pbOrg.ManagerItem
		if !o.ManagerID.IsZero() {
			val, ok := mManager[o.ManagerID]
			if !ok {
				err = errors.New("expected to get manager from map, but nothing found o.ManagerID=" + o.ManagerID.Hex())
				return
			}

			managerItem = &pbOrg.ManagerItem{
				Id:   val.ID.Hex(),
				Slug: val.Slug,
				Name: val.Name,
			}
		}

		var okvedOsn *pbOrg.OkvedFullItem
		if !o.OkvedOsnID.IsZero() {
			val, ok := mOkved[o.OkvedOsnID]
			if !ok {
				err = errors.New("expected to get okvedOsn from map, but nothing found o.OkvedOsnID=" + o.OkvedOsnID.Hex())
				return
			}

			okvedOsn = &pbOrg.OkvedFullItem{
				Id:           val.ID.Hex(),
				Slug:         val.Slug,
				Name:         val.Name,
				Code:         val.Code,
				CodeWithName: val.CodeWithName,
				Kind:         pbOrg.OkvedKind(val.Kind),
			}
		}

		var okvedDop []*pbOrg.OkvedItem
		for _, id := range o.OkvedDopIDs {
			val, ok := mOkved[id]
			if !ok {
				err = errors.New("expected to get okvedOsn from map, but nothing found o.OkvedDopIDs=" + id.Hex())
				return
			}

			okvedDop = append(okvedDop, &pbOrg.OkvedItem{
				Id:   val.ID.Hex(),
				Slug: val.Slug,
				Name: val.Name,
				Kind: pbOrg.OkvedKind(val.Kind),
			})
		}

		var metros []*pbOrg.MetroFullItemWithDistance
		for _, m := range o.Metros {
			if !m.ID.IsZero() {
				val, ok := mMetro[m.ID]
				if !ok {
					err = errors.New("expected to get metro from map, but nothing found m.ID=" + m.ID.Hex())
					return
				}

				ar, ok := mArea[val.AreaID]
				if !ok {
					err = errors.New("expected to get area from map, but nothing found val.AreaID=" + val.AreaID.Hex())
					return
				}

				metros = append(metros, &pbOrg.MetroFullItemWithDistance{
					Id:       val.ID.Hex(),
					Slug:     val.Slug,
					Name:     val.Name,
					Line:     val.Line,
					Distance: m.Distance,
					Area: &pbOrg.AreaFullItem{
						Id:       ar.ID.Hex(),
						Slug:     ar.Slug,
						Name:     ar.Name,
						FiasId:   ar.FiasID,
						KladrId:  ar.KladrID,
						Type:     ar.Type,
						TypeFull: ar.TypeFull,
					},
				})
			}
		}

		res.Main = &pbOrg.Main{
			Id:               o.ID.Hex(),
			Slug:             o.Slug,
			Name:             o.Name,
			Inn:              float64(o.INN),
			Kpp:              float64(o.KPP),
			Ogrn:             float64(o.OGRN),
			Kind:             pbOrg.OrgKind(o.Kind),
			Manager:          managerItem,
			Area:             areaFullItem,
			Location:         locationItem,
			Okved:            okvedOsn,
			StatusKind:       pbOrg.StatusKind(o.StatusKind),
			OkvedDop:         okvedDop,
			EmployeeCount:    o.EmployeeCount,
			Metros:           metros,
			NameFullWithOpf:  o.NameFullWithOPF,
			NameShortWithOpf: o.NameShortWithOPF,
			OpfCode:          float64(o.OPFCode),
			OpfFull:          o.OPFFull,
			OpfShort:         o.OPFShort,
			OpfKind:          pbOrg.OpfKind(o.OPFKind),
			OgrnDate:         o.OGRNDate.String(),
			Okato:            float64(o.OKATO),
			Oktmo:            float64(o.OKTMO),
			Okpo:             float64(o.OKPO),
			Okogu:            float64(o.OKOGU),
			Okfs:             float64(o.OKFS),
			RegistrationDate: toNotZeroTime(o.RegistrationDate),
			LiquidationDate:  toNotZeroTime(o.LiquidationDate),
			UpdatedAt:        toNotZeroTime(o.UpdatedAt),
		}
		res.Related = related
	}
	return
}

func toNotZeroTime(in time.Time) (out string) {
	if in.After(time.Unix(0, 0).UTC()) {
		out = in.String()
	}
	return
}
