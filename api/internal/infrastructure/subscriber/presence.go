package subscriber

import (
	"context"
	"log"

	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/icchon/matcha/api/internal/domain/client"
)

const presenceChannel = "presence_incoming"

type presenceSubscriber struct {
	rdb     *redis.Client
	channel string
}

func NewPresenceSubscriber(rdb *redis.Client) *presenceSubscriber {
	return &presenceSubscriber{
		rdb:     rdb,
		channel: presenceChannel,
	}
}

var _ client.Subscriber = (*presenceSubscriber)(nil)

func (s *presenceSubscriber) SubscribeChannel(ctx context.Context, handler func(ctx context.Context, payload interface{}) error) error {
	pubsub := s.rdb.Subscribe(ctx, string(s.channel))
	ch := pubsub.Channel()

	go func() {
		defer pubsub.Close()
		for {
			select {
			case <-ctx.Done():
				log.Printf("Context cancelled. Stopping subscription for channel: %s", s.channel)
				return

			case msg, ok := <-ch:
				if !ok {
					log.Printf("Redis channel closed for %s.", s.channel)
					return
				}

				var payload client.PresencePayload
				err := json.Unmarshal([]byte(msg.Payload), &payload)
				if err != nil {
					log.Printf("Error unmarshaling message from channel %s: %v", s.channel, err)
					continue
				}

				if err := handler(ctx, &payload); err != nil {
					log.Printf("Error handling message from channel %s: %v", s.channel, err)
				}
			}
		}
	}()
	return nil
}
