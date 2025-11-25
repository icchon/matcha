package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestMessageRepository_GetLatest(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	r := NewMessageRepository(db)

	t.Run("Get latest message", func(t *testing.T) {
		expectedSQL := `SELECT \* FROM messages WHERE \(sender_id = \$1 AND recipient_id = \$2\) OR \(sender_id = \$2 AND recipient_id = \$1\) ORDER BY sent_at DESC LIMIT 1`

		rows := sqlmock.NewRows([]string{"id", "sender_id", "recipient_id", "content", "sent_at", "is_read"}).
			AddRow(1, userID1, userID2, "hello", time.Now(), false)

		mock.ExpectQuery(expectedSQL).
			WithArgs(userID1, userID2).
			WillReturnRows(rows)

		msg, err := r.GetLatest(context.Background(), userID1, userID2)

		assert.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, "hello", msg.Content)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMessageRepository_Query(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	r := NewMessageRepository(db)

	t.Run("Get messages between two users with pagination", func(t *testing.T) {
		query := &repo.MessageQuery{
			SenderID:    &userID1,
			RecipientID: &userID2,
			Limit:       func(i int) *int { return &i }(10),
			Offset:      func(i int) *int { return &i }(0),
		}
		expectedSQL := `SELECT \* FROM messages WHERE 1=1 AND \(\(sender_id = \$1 AND recipient_id = \$2\) OR \(sender_id = \$2 AND recipient_id = \$1\)\) ORDER BY sent_at DESC LIMIT \$3 OFFSET \$4`

		rows := sqlmock.NewRows([]string{"id", "sender_id", "recipient_id", "content", "sent_at", "is_read"}).
			AddRow(1, userID1, userID2, "hello", time.Now(), false)

		mock.ExpectQuery(expectedSQL).
			WithArgs(userID1, userID2, 10, 0).
			WillReturnRows(rows)

		messages, err := r.Query(context.Background(), query)

		assert.NoError(t, err)
		assert.Len(t, messages, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
