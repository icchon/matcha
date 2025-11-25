package repo

import (
	"context"
	"github.com/icchon/matcha/api/internal/domain/entity"
)

type TagRepository interface {
	GetAll(ctx context.Context) ([]*entity.Tag, error)
}
