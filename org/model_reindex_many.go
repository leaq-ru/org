package org

import (
	"context"
	"github.com/gosimple/slug"
	"github.com/nnqq/scr-org/dadata"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

type Upsert struct {
	AreaID        primitive.ObjectID
	LocationID    primitive.ObjectID
	ManagerID     primitive.ObjectID
	OkvedOsnID    primitive.ObjectID
	OkvedDopIDs   []primitive.ObjectID
	Metros        []Metro
	Sugg          dadata.Suggestion
	EmployeeCount uint32
}

func (m Model) ReindexMany(
	ctx context.Context,
	vals []Upsert,
) (
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if len(vals) == 0 {
		return
	}

	var wm []mongo.WriteModel
	for _, v := range vals {
		opfCode, e := toUint64(v.Sugg.Data.Opf.Code)
		if e != nil {
			err = e
			return
		}
		inn, e := toUint64(v.Sugg.Data.Inn)
		if e != nil {
			err = e
			return
		}
		kpp, e := toUint64(v.Sugg.Data.Kpp)
		if e != nil {
			err = e
			return
		}
		ogrn, e := toUint64(v.Sugg.Data.Ogrn)
		if e != nil {
			err = e
			return
		}
		okato, e := toUint64(v.Sugg.Data.Okato)
		if e != nil {
			err = e
			return
		}
		oktmo, e := toUint64(v.Sugg.Data.Oktmo)
		if e != nil {
			err = e
			return
		}
		okpo, e := toUint64(v.Sugg.Data.Okpo)
		if e != nil {
			err = e
			return
		}
		okogu, e := toUint64(v.Sugg.Data.Okogu)
		if e != nil {
			err = e
			return
		}
		okfs, e := toUint64(v.Sugg.Data.Okfs)
		if e != nil {
			err = e
			return
		}

		uo := mongo.NewUpdateOneModel()
		uo.SetFilter(org{
			Slug: makeSlug(v.Sugg),
		})
		uo.SetUpdate(bson.M{
			"$set": org{
				DaDataID:      v.Sugg.Data.Hid,
				AreaID:        v.AreaID,
				LocationID:    v.LocationID,
				ManagerID:     v.ManagerID,
				ManagerPost:   v.Sugg.Data.Management.Post,
				EmployeeCount: v.EmployeeCount,
				OkvedOsnID:    v.OkvedOsnID,
				OkvedDopIDs:   v.OkvedDopIDs,
				Metros:        v.Metros,
				Name: strings.Join([]string{
					v.Sugg.Data.Opf.Short,
					v.Sugg.Data.Name.Full,
				}, " "),
				OPFCode:          opfCode,
				OPFFull:          v.Sugg.Data.Opf.Full,
				OPFShort:         v.Sugg.Data.Opf.Short,
				OPFKind:          toOPFKind(v.Sugg.Data.Opf.Type),
				Kind:             toKind(v.Sugg.Data.Type),
				BranchKind:       toBranchKind(v.Sugg.Data.BranchType),
				BranchCount:      uint32(v.Sugg.Data.BranchCount),
				INN:              inn,
				KPP:              kpp,
				OGRN:             ogrn,
				OGRNDate:         msTsToTime(v.Sugg.Data.OgrnDate),
				OKATO:            okato,
				OKTMO:            oktmo,
				OKPO:             okpo,
				OKOGU:            okogu,
				OKFS:             okfs,
				StatusKind:       toStatusKind(v.Sugg.Data.State.Status),
				RegistrationDate: msTsToTime(v.Sugg.Data.State.RegistrationDate),
				LiquidationDate:  msTsToTime(v.Sugg.Data.State.LiquidationDate),
				UpdatedAt:        time.Now().UTC(),
			},
		})
		uo.SetUpsert(true)
		wm = append(wm, uo)
	}

	_, err = m.coll.BulkWrite(ctx, wm, options.BulkWrite().SetOrdered(false))
	return
}

func toUint64(in string) (out uint64, err error) {
	if in == "" {
		return
	}

	i, err := strconv.Atoi(in)
	out = uint64(i)
	return
}

func toOPFKind(in string) opfKind {
	switch in {
	case "99":
		return opfKind_y1999
	case "2012":
		return opfKind_y2012
	case "2014":
		return opfKind_y2014
	default:
		return 0
	}
}

func toKind(in string) kind {
	switch in {
	case "LEGAL":
		return kind_legal
	case "INDIVIDUAL":
		return kind_individual
	default:
		return 0
	}
}

func toBranchKind(in string) branchKind {
	switch in {
	case "MAIN":
		return branchKind_main
	case "BRANCH":
		return branchKind_branch
	default:
		return 0
	}
}

func toStatusKind(in string) statusKind {
	switch in {
	case "ACTIVE":
		return statusKind_active
	case "LIQUIDATING":
		return statusKind_liquidating
	case "LIQUIDATED":
		return statusKind_liquidated
	case "BANKRUPT":
		return statusKind_bankrupt
	case "REORGANIZING":
		return statusKind_reorganizing
	default:
		return 0
	}
}

func msTsToTime(in int64) time.Time {
	return time.Unix(in/1000, 0)
}

func makeSlug(sugg dadata.Suggestion) string {
	return slug.Make(strings.Join([]string{
		sugg.Data.Opf.Short,
		sugg.Data.Name.Full,
		sugg.Data.Inn,
	}, " "))
}
