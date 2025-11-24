package publisher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/api/internal/domain/client"
)

const notificationChannel string = "notification_channel"

type notificationPublisher struct {
	rdb     *redis.Client
	channel string
}

var _ client.Publisher = (*notificationPublisher)(nil)

func NewNotificationPublisher(rdb *redis.Client) *notificationPublisher {
	return &notificationPublisher{
		rdb:     rdb,
		channel: notificationChannel,
	}
}

func (p *notificationPublisher) Publish(ctx context.Context, data interface{}) error {
	return p.rdb.Publish(ctx, p.channel, data).Err()
}
