package publisher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/api/internal/domain/client"
)

const chatCannnel string = "chat_channel"

type chatPublisher struct {
	rdb     *redis.Client
	channel string
}

var _ client.Publisher = (*chatPublisher)(nil)

func NewChatPublisher(rdb *redis.Client) *chatPublisher {
	return &chatPublisher{
		rdb:     rdb,
		channel: chatCannnel,
	}
}

func (p *chatPublisher) Publish(ctx context.Context, data interface{}) error {
	p.rdb.Publish(ctx, p.channel, data)
	return nil
}
