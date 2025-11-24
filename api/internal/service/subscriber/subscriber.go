package subscriber

import (
	"context"

	"github.com/icchon/matcha/api/internal/domain/client"
	"github.com/icchon/matcha/api/internal/domain/service"
)

type SubscriberHandler interface {
	ChatSubscHandler(ctx context.Context, payload *client.MessagePayload) error
	PresenceSubscHandler(ctx context.Context, payload *client.PresencePayload) error
	ReadSubscHandler(ctx context.Context, payload *client.ReadPayload) error
}

type subscriberService struct {
	chatSub     client.Subscriber
	presenseSub client.Subscriber
	readSub     client.Subscriber

	subscHandler SubscriberHandler
}

var _ service.SubscriberService = (*subscriberService)(nil)

func NewSubscriberService(
	chatSub client.Subscriber,
	presenseSub client.Subscriber,
	readSub client.Subscriber,
	subscHandler SubscriberHandler,
) *subscriberService {
	return &subscriberService{
		chatSub:      chatSub,
		presenseSub:  presenseSub,
		readSub:      readSub,
		subscHandler: subscHandler,
	}
}

func (s *subscriberService) Initialize(ctx context.Context) error {
	if err := s.chatSub.SubscribeChannel(ctx, func(ctx context.Context, data interface{}) error {
		s.subscHandler.ChatSubscHandler(ctx, data.(*client.MessagePayload))
		return nil
	}); err != nil {
		return err
	}
	if err := s.presenseSub.SubscribeChannel(ctx, func(ctx context.Context, data interface{}) error {
		s.subscHandler.PresenceSubscHandler(ctx, data.(*client.PresencePayload))
		return nil
	}); err != nil {
		return err
	}
	if err := s.readSub.SubscribeChannel(ctx, func(ctx context.Context, data interface{}) error {
		s.subscHandler.ReadSubscHandler(ctx, data.(*client.ReadPayload))
		return nil
	}); err != nil {
		return err
	}
	return nil
}
