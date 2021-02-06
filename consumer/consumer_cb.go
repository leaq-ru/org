package consumer

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-org-producer/protocol"
	"github.com/nnqq/scr-org/area"
	"github.com/nnqq/scr-org/dadata"
	"github.com/nnqq/scr-org/metro"
	"github.com/nnqq/scr-org/okved"
	"github.com/nnqq/scr-org/org"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
	"strconv"
	"time"
)

func (c Consumer) cb(rawMsg *stan.Msg) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		var msg protocol.OrgMessage
		err := json.Unmarshal(rawMsg.Data, &msg)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		ack := func() {
			e := rawMsg.Ack()
			if e != nil {
				c.logger.Error().Err(e).Send()
			}
		}

		suggs, err := c.dadataClient.GetByINN(ctx, msg.INN)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}
		if len(suggs) == 0 {
			ack()
			return
		}

		ups := make([]org.Upsert, len(suggs))
		var eg errgroup.Group
		for _i, _sugg := range suggs {
			i := _i
			sugg := _sugg

			ups[i].Sugg = sugg
			if msg.EmployeeCount != "" {
				ec, e := strconv.Atoi(msg.EmployeeCount)
				if e != nil {
					c.logger.Error().Err(e).Send()
					return
				}
				ups[i].EmployeeCount = uint32(ec)
			}

			eg.Go(func() (e error) {
				ups[i].LocationID, e = c.locationModel.Find(ctx, sugg.Data.Address.UnrestrictedValue)
				return
			})

			eg.Go(func() (e error) {
				ups[i].ManagerID, e = c.managerModel.Find(ctx, sugg.Data.Management.Name)
				return
			})

			eg.Go(func() error {
				var ou []okved.Upsert
				if msg.OkvedOsn.Code != "" && msg.OkvedOsn.Name != "" && msg.OkvedOsn.Ver != "" {
					ou = append(ou, okved.Upsert{
						Code: msg.OkvedOsn.Code,
						Name: msg.OkvedOsn.Name,
						Kind: msg.OkvedOsn.Ver,
					})
				}

				ids, e := c.okvedModel.FindMany(ctx, ou)
				if e != nil {
					return e
				}

				ups[i].OkvedOsnID = ids[0]
				return nil
			})

			eg.Go(func() (e error) {
				var ou []okved.Upsert
				for _, od := range msg.OkvedDop {
					ou = append(ou, okved.Upsert{
						Code: od.Code,
						Name: od.Name,
						Kind: od.Ver,
					})
				}

				ups[i].OkvedDopIDs, e = c.okvedModel.FindMany(ctx, ou)
				return
			})

			eg.Go(func() error {
				areaID, e := c.getAreaID(ctx, sugg.Data.Address.Data)
				if e != nil {
					return e
				}
				ups[i].AreaID = areaID

				var fmReq []metro.FindManyReqItem
				for _, m := range sugg.Data.Address.Data.Metro {
					fmReq = append(fmReq, metro.FindManyReqItem{
						Name:     m.Name,
						Line:     m.Line,
						Distance: float32(m.Distance),
					})
				}

				metrosRaw, e := c.metroModel.FindMany(ctx, areaID, fmReq)
				if e != nil {
					return e
				}

				for _, mr := range metrosRaw {
					ups[i].Metros = append(ups[i].Metros, org.Metro{
						ID:       mr.ID,
						Distance: mr.Distance,
					})
				}
				return nil
			})
		}
		err = eg.Wait()
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		err = c.orgModel.ReindexMany(ctx, ups)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		ack()
	}()
}

func (c Consumer) getAreaID(ctx context.Context, ad dadata.AddressData) (id primitive.ObjectID, err error) {
	var areaKind area.Kind
	if ad.City != "" {
		areaKind = area.Kind_city
	} else if ad.Settlement != "" {
		areaKind = area.Kind_settlement
	} else {
		return
	}

	return c.areaModel.Find(
		ctx,
		notEmpty(ad.City, ad.Settlement),
		notEmpty(ad.CityFiasID, ad.SettlementFiasID),
		notEmpty(ad.CityKladrID, ad.SettlementKladrID),
		notEmpty(ad.CityType, ad.SettlementType),
		notEmpty(ad.CityTypeFull, ad.SettlementTypeFull),
		areaKind,
	)
}

func notEmpty(str1, str2 string) string {
	if str1 != "" {
		return str1
	}
	return str2
}
