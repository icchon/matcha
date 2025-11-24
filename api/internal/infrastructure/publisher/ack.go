package publisher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/api/internal/domain/client"
)

const ackChannel string = "ack_channel"

type ackPublisher struct {
	rdb     *redis.Client
	channel string
}

var _ client.Publisher = (*ackPublisher)(nil)

func NewAckPublisher(rdb *redis.Client) *ackPublisher {
	return &ackPublisher{
		rdb:     rdb,
		channel: ackChannel,
	}
}

func (p *ackPublisher) Publish(ctx context.Context, data interface{}) error {
	return p.rdb.Publish(ctx, p.channel, data).Err()
}
