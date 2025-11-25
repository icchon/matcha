package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestConnectionRepository_Query(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	r := NewConnectionRepository(db)

	t.Run("Find connections for a user", func(t *testing.T) {
		query := &repo.ConnectionQuery{User1ID: &userID1}
		expectedSQL := `SELECT \* FROM connections WHERE 1=1 AND \(user1_id = \$1 OR user2_id = \$1\)`

		rows := sqlmock.NewRows([]string{"user1_id", "user2_id"}).
			AddRow(userID1, userID2).
			AddRow(userID3, userID1)

		mock.ExpectQuery(expectedSQL).
			WithArgs(userID1).
			WillReturnRows(rows)

		connections, err := r.Query(context.Background(), query)

		assert.NoError(t, err)
		assert.Len(t, connections, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
