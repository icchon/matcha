package publisher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/api/internal/domain/client"
)

const notificationCannnel string = "notification_channel"

type notificationPublisher struct {
	rdb     *redis.Client
	channel string
}

var _ client.Publisher = (*ackPublisher)(nil)

func NewNotificationPublisher(rdb *redis.Client) *notificationPublisher {
	return &notificationPublisher{
		rdb:     rdb,
		channel: notificationCannnel,
	}
}

func (p *notificationPublisher) Publish(ctx context.Context, data interface{}) error {
	p.rdb.Publish(ctx, p.channel, data)
	return nil
}
