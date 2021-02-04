package consumer

import (
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/nnqq/scr-org-producer/protocol"
	"time"
)

func (c Consumer) cb(rawMsg *stan.Msg) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		ack := func() {
			e := rawMsg.Ack()
			if e != nil {
				c.logger.Error().Err(e).Send()
			}
		}

		var msg protocol.OrgMessage
		err := json.Unmarshal(rawMsg.Data, &msg)
		if err != nil {
			c.logger.Error().Err(err).Send()
			return
		}

	}()
}
