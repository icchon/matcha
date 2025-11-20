package helper

import (
	"context"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/presentation/middleware"
)

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(middleware.UserIDContextKey).(uuid.UUID)
	return userID, ok
}
