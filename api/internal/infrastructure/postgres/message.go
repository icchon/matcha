package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type messageRepository struct {
	db DBTX
}

func NewMessageRepository(db DBTX) repo.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *entity.Message) error {
	query := `
		INSERT INTO messages (sender_id, recipient_id, content)
		VALUES (:sender_id, :recipient_id, :content)
		RETURNING *
	`
	return r.db.QueryRowxContext(ctx, query, message).StructScan(message)
}

func (r *messageRepository) Update(ctx context.Context, message *entity.Message) error {
	query := `
		UPDATE messages SET
			content = :content,
			is_read = :is_read
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, message)
	return err
}

func (r *messageRepository) Delete(ctx context.Context, messageID int64) error {
	query := "DELETE FROM messages WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, messageID)
	return err
}

func (r *messageRepository) Find(ctx context.Context, messageID int64) (*entity.Message, error) {
	var message entity.Message
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

func (r *messageRepository) Query(ctx context.Context, q *repo.MessageQuery) ([]*entity.Message, error) {
	query := "SELECT * FROM messages WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.ID != nil {
		query += fmt.Sprintf(" AND id = $%d", argCount)
		args = append(args, *q.ID)
		argCount++
	}
	if q.SenderID != nil {
		query += fmt.Sprintf(" AND sender_id = $%d", argCount)
		args = append(args, *q.SenderID)
		argCount++
	}
	if q.RecipientID != nil {
		query += fmt.Sprintf(" AND recipient_id = $%d", argCount)
		args = append(args, *q.RecipientID)
		argCount++
	}
	if q.Content != nil {
		query += fmt.Sprintf(" AND content LIKE $%d", argCount)
		args = append(args, "%"+*q.Content+"%")
		argCount++
	}
	if q.SentAt != nil {
		query += fmt.Sprintf(" AND sent_at = $%d", argCount)
		args = append(args, *q.SentAt)
		argCount++
	}
	if q.IsRead != nil {
		query += fmt.Sprintf(" AND is_read = $%d", argCount)
		args = append(args, *q.IsRead)
		argCount++
	}

	var messages []*entity.Message
	if err := r.db.SelectContext(ctx, &messages, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return messages, nil
}
