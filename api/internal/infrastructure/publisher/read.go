package publisher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/api/internal/domain/client"
)

const readChannel string = "read_outgoing"

type readPublisher struct {
	rdb     *redis.Client
	channel string
}

var _ client.Publisher = (*readPublisher)(nil)

func NewReadPublisher(rdb *redis.Client) *readPublisher {
	return &readPublisher{
		rdb:     rdb,
		channel: readChannel,
	}
}

func (p *readPublisher) Publish(ctx context.Context, data interface{}) error {
	return p.rdb.Publish(ctx, p.channel, data).Err()
}
