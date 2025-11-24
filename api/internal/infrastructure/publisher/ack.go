package publisher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/api/internal/domain/client"
)

const ackCannnel string = "ack_channel"

type ackPublisher struct {
	rdb     *redis.Client
	channel string
}

var _ client.Publisher = (*ackPublisher)(nil)

func NewAckPublisher(rdb *redis.Client) *ackPublisher {
	return &ackPublisher{
		rdb:     rdb,
		channel: ackCannnel,
	}
}

func (p *ackPublisher) Publish(ctx context.Context, data interface{}) error {
	p.rdb.Publish(ctx, p.channel, data)
	return nil
}
