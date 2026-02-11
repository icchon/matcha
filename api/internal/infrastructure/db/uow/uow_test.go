package uow

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUnitOfWork_Do(t *testing.T) {
	fnErr := errors.New("fn error")

	testCases := []struct {
		name        string
		fnError     error
		expectation func(mock sqlmock.Sqlmock)
		expectedErr error
	}{
		{
			name:    "Success - fn succeeds and commit is called",
			fnError: nil,
			expectation: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name:    "fn returns error - Do returns fn error, not rollback error",
			fnError: fnErr,
			expectation: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			expectedErr: fnErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err, "Failed to create sqlmock. Check go-sqlmock setup.")
			defer db.Close()

			tc.expectation(mock)

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			uow := NewUnitOfWork(sqlxDB)

			err = uow.Do(context.Background(), func(m repo.RepositoryManager) error {
				return tc.fnError
			})

			assert.Equal(t, tc.expectedErr, err,
				"Do() should return the fn error when fn fails, not the rollback error. Check that uow.go returns err instead of tx.Rollback().")
			assert.NoError(t, mock.ExpectationsWereMet(),
				"SQL expectations not met. Verify that Rollback is called on fn error, and Commit on fn success.")
		})
	}
}
