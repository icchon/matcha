package publisher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/api/internal/domain/client"
)

const chatChannel string = "chat_outgoing"

type chatPublisher struct {
	rdb     *redis.Client
	channel string
}

var _ client.Publisher = (*chatPublisher)(nil)

func NewChatPublisher(rdb *redis.Client) *chatPublisher {
	return &chatPublisher{
		rdb:     rdb,
		channel: chatChannel,
	}
}

func (p *chatPublisher) Publish(ctx context.Context, data interface{}) error {
	return p.rdb.Publish(ctx, p.channel, data).Err()
}
