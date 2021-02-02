package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-exporter/cached_export"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func (c Consumer) cb(rawMsg *stan.Msg) {
	go func() {
		deadline := 10 * time.Hour
		commitStateInterval := 10 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), deadline)
		defer cancel()

		ack := func() {
			e := rawMsg.Ack()
			if e != nil {
				c.logger.Error().Err(e).Send()
			}
		}

		var msg message
		err := json.Unmarshal(rawMsg.Data, &msg)
		if err != nil {
			c.logger.Error().Err(err).Msg("got malformed msg, just ack")
			ack()
			return
		}

		setFail := func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := c.fileModel.SetFail(ctx, msg.ID)
			if err != nil {
				c.logger.Error().Err(err).Send()
			}
		}

		if rawMsg.Timestamp < time.Now().UTC().Add(-deadline).UnixNano() {
			setFail()
			ack()
			return
		}

		processing, err := c.processingExportModel.IsProcessing(ctx, msg.ID)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		if processing {
			return
		}

		unsetProcessing := func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := c.processingExportModel.UnsetProcessing(ctx, msg.ID)
			if err != nil {
				c.logger.Error().Err(err).Send()
			}
		}

		go func() {
			for ctx.Err() == nil {
				err := c.processingExportModel.SetProcessing(ctx, msg.ID)
				if err != nil {
					c.logger.Error().Err(err).Send()
				}
				time.Sleep(commitStateInterval)
			}
		}()
		defer unsetProcessing()

		reqComp := &parser.GetV2Request{
			CityIds:            msg.Query.GetCityIds(),
			CategoryIds:        msg.Query.GetCategoryIds(),
			HasEmail:           msg.Query.GetHasEmail(),
			HasPhone:           msg.Query.GetHasPhone(),
			HasOnline:          msg.Query.GetHasOnline(),
			HasInn:             msg.Query.GetHasInn(),
			HasKpp:             msg.Query.GetHasKpp(),
			HasOgrn:            msg.Query.GetHasOgrn(),
			HasAppStore:        msg.Query.GetHasAppStore(),
			HasGooglePlay:      msg.Query.GetHasGooglePlay(),
			HasVk:              msg.Query.GetHasVk(),
			VkMembersCount:     msg.Query.GetVkMembersCount(),
			HasInstagram:       msg.Query.GetHasInstagram(),
			HasTwitter:         msg.Query.GetHasTwitter(),
			HasYoutube:         msg.Query.GetHasYoutube(),
			HasFacebook:        msg.Query.GetHasFacebook(),
			TechnologyIds:      msg.Query.GetTechnologyIds(),
			TechnologyFindRule: msg.Query.GetTechnologyFindRule(),
			DnsIds:             msg.Query.GetDnsIds(),
		}

		masterJob, err := c.fileModel.GetMasterJob(ctx, reqComp, msg.ID)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		if !masterJob.IsZero() {
			err := c.fileModel.WatchJob(ctx, masterJob, msg.ID)
			if err != nil {
				c.logger.Error().Err(err).Send()
				return
			}

			ack()
			return
		}

		defer func() {
			err := c.rowModel.Flush(ctx)
			if err != nil {
				c.logger.Error().Err(err).Send()
			}
		}()

		fromCompanyID, err := c.fileModel.GetFromCompanyID(ctx, msg.ID)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}
		saveFromCompanyID := func() {
			if fromCompanyID == "" {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			e := c.fileModel.SetFromCompanyID(ctx, msg.ID, fromCompanyID)
			if e != nil {
				c.logger.Error().Err(e).Send()
			}
		}

		var loopDone bool
		go func() {
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

			select {
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					unsetProcessing()
					setFail()
					ack()
				}
			case <-signals:
				var wg sync.WaitGroup
				wg.Add(4)
				go func() {
					defer wg.Done()
					unsetProcessing()
				}()
				go func() {
					defer wg.Done()
					saveFromCompanyID()
				}()
				go func() {
					defer wg.Done()
					ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
					defer cancel()
					err := c.rowModel.Flush(ctx)
					if err != nil {
						c.logger.Error().Err(err).Send()
					}
				}()
				go func() {
					defer wg.Done()
					for {
						if loopDone {
							return
						}
						time.Sleep(time.Second)
					}
				}()
				wg.Wait()
			}
		}()

		s3URL, err := c.cachedExportModel.Get(ctx, reqComp)
		if err != nil {
			if errors.Is(err, cached_export.ErrNoFound) {
				err = nil

				go func() {
					resCount, err := c.companyClient.GetCount(ctx, reqComp)
					if err != nil {
						c.logger.Error().Err(err).Send()
						return
					}

					err = c.fileModel.SetMasterJobInProgress(ctx, msg.ID, reqComp, resCount.GetCount())
					if err != nil {
						c.logger.Error().Err(err).Send()
						return
					}
				}()

				compStream, err := c.companyClient.GetFull(ctx, &parser.GetFullRequest{
					Query:  reqComp,
					FromId: fromCompanyID,
				})
				if err != nil {
					c.logger.Error().Err(err).Send()
					return
				}

				var (
					mu           sync.Mutex
					currentCount uint32
				)
				go func() {
					for {
						select {
						case <-ctx.Done():
							return
						default:
							time.Sleep(commitStateInterval)

							var eg errgroup.Group
							eg.Go(func() (e error) {
								saveFromCompanyID()
								return
							})
							eg.Go(func() (e error) {
								mu.Lock()
								delta := currentCount
								currentCount = 0
								mu.Unlock()

								if delta == 0 {
									return
								}
								e = c.fileModel.IncCurrentCount(ctx, msg.ID, delta)
								return
							})
							err = eg.Wait()
							if err != nil {
								c.logger.Error().Err(err).Send()
							}
						}
					}
				}()

				for {
					if c.state.drain {
						loopDone = true
						return
					}

					comp, err := compStream.Recv()
					if err != nil {
						if errors.Is(err, io.EOF) {
							break
						}

						c.logger.Error().Err(err).Send()
						return
					}

					err = c.rowModel.Add(ctx, msg.ID, comp)
					if err != nil {
						c.logger.Error().Err(err).Send()
						return
					}

					mu.Lock()
					currentCount += 1
					mu.Unlock()
					fromCompanyID = comp.GetId()
				}

				err = c.rowModel.Flush(ctx)
				if err != nil {
					c.logger.Error().Err(err).Send()
					return
				}

				csvPath, err := c.rowModel.PipeMongoToLocalFile(ctx, msg.ID)
				if err != nil {
					c.logger.Error().Err(err).Send()
					return
				}

				s3URL, err = c.exporterBucket.Put(ctx, csvPath, true)
				if err != nil {
					c.logger.Error().Err(err).Send()
					return
				}

				err = c.cachedExportModel.Set(ctx, reqComp, s3URL)
				if err != nil {
					c.logger.Error().Err(err).Send()
					return
				}
			} else {
				c.logger.Error().Err(err).Send()
				return
			}
		}

		err = c.fileModel.SetSuccess(ctx, msg.ID, s3URL)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

		ack()
		return
	}()
}
