package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type messageRepository struct {
	db DBTX
}

func NewMessageRepository(db DBTX) repo.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) MarkAsRead(ctx context.Context, messageID int64) error {
	query := `
		UPDATE messages SET
			is_read = TRUE
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, messageID)
	return err
}

func (r *messageRepository) Create(ctx context.Context, message *entity.Message) error {
	query := `
		INSERT INTO messages (sender_id, recipient_id, content)
		VALUES (:sender_id, :recipient_id, :content)
		RETURNING *
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRowxContext(ctx, message).StructScan(message)
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

	if q.SenderID != nil && q.RecipientID != nil {
		query += fmt.Sprintf(" AND ((sender_id = $%d AND recipient_id = $%d) OR (sender_id = $%d AND recipient_id = $%d))", argCount, argCount+1, argCount+1, argCount)
		args = append(args, *q.SenderID, *q.RecipientID)
		argCount += 2
	} else if q.SenderID != nil {
		query += fmt.Sprintf(" AND sender_id = $%d", argCount)
		args = append(args, *q.SenderID)
		argCount++
	} else if q.RecipientID != nil {
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

	query += " ORDER BY sent_at DESC"

	if q.Limit != nil {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, *q.Limit)
		argCount++
	}
	if q.Offset != nil {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, *q.Offset)
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

func (r *messageRepository) GetLatest(ctx context.Context, userID1, userID2 uuid.UUID) (*entity.Message, error) {
	var message entity.Message
	query := `
		SELECT *
		FROM messages
		WHERE (sender_id = $1 AND recipient_id = $2) OR (sender_id = $2 AND recipient_id = $1)
		ORDER BY sent_at DESC
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &message, query, userID1, userID2)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No message found, which is a valid scenario
		}
		return nil, err
	}
	return &message, nil
}
