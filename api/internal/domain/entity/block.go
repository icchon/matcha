package entity

import (
	"github.com/google/uuid"
)

type Block struct {
	BlockerID uuid.UUID `db:"blocker_id"`
	BlockedID uuid.UUID `db:"blocked_id"`
}
