package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type messageRepository struct {
	db DBTX
}

func NewMessageRepository(db DBTX) domain.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Save(ctx context.Context, message *domain.Message) error {
	query := `
		INSERT INTO messages (id, sender_id, recipient_id, content, sent_at, is_read)
		VALUES (:id, :sender_id, :recipient_id, :content, :sent_at, :is_read)
		ON CONFLICT (id) DO UPDATE SET
			content = :content,
			is_read = :is_read
	`
	_, err := r.db.NamedExecContext(ctx, query, message)
	return err
}

func (r *messageRepository) Delete(ctx context.Context, messageID int64) error {
	query := "DELETE FROM messages WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, messageID)
	return err
}

func (r *messageRepository) Find(ctx context.Context, messageID int64) (*domain.Message, error) {
	var message domain.Message
	query := "SELECT * FROM messages WHERE id = $1"
	err := r.db.GetContext(ctx, &message, query, messageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) GetBySenderID(ctx context.Context, senderID uuid.UUID) ([]domain.Message, error) {
	var messages []domain.Message
	query := "SELECT * FROM messages WHERE sender_id = $1"
	err := r.db.SelectContext(ctx, &messages, query, senderID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
