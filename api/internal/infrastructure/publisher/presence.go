package publisher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/api/internal/domain/client"
)

const presenceCannnel string = "presence_channel"

type presencePublisher struct {
	rdb     *redis.Client
	channel string
}

var _ client.Publisher = (*presencePublisher)(nil)

func NewPresencePublisher(rdb *redis.Client) *presencePublisher {
	return &presencePublisher{
		rdb:     rdb,
		channel: presenceCannnel,
	}
}

func (p *presencePublisher) Publish(ctx context.Context, data interface{}) error {
	p.rdb.Publish(ctx, p.channel, data)
	return nil
}
