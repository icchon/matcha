package service

import (
	"context"
)

type SubscriberService interface {
	Initialize(ctx context.Context) error
}
