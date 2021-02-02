package stan

import (
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"strings"
)

func NewConn(serviceName, clusterID, natsURL string) (stan.Conn, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return stan.Connect(
		clusterID,
		strings.Join([]string{
			serviceName,
			u.String(),
		}, "-"),
		stan.NatsURL(natsURL),
	)
}
