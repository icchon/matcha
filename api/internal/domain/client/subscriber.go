package client

import (
	"context"
)

type Subscriber interface {
	SubscribeChannel(ctx context.Context, handler func(ctx context.Context, data interface{}) error) error
}
