package consumer

import (
	"github.com/nats-io/stan.go"
	"time"
)

const orgSubjectName = "org"

func (c Consumer) Subscribe() (err error) {
	c.state.sub, err = c.stanConn.QueueSubscribe(
		orgSubjectName,
		c.serviceName,
		c.cb,
		stan.DurableName(orgSubjectName),
		stan.SetManualAckMode(),
		stan.AckWait(30*time.Second),
		stan.MaxInflight(1),
	)
	if err != nil {
		return
	}

	c.state.subscribeCalledOK = true
	return
}
