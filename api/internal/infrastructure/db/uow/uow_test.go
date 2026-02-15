package uow

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/icchon/matcha/api/internal/domain/repo" // Add this import
)

func TestUnitOfWork_Do(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	uow := NewUnitOfWork(sqlxDB)

	t.Run("successful transaction commit", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		err := uow.Do(context.Background(), func(rm repo.RepositoryManager) error {
			return nil
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("transaction rollback on function error", func(t *testing.T) {
		testErr := errors.New("function error")
		mock.ExpectBegin()
		mock.ExpectRollback()

		err := uow.Do(context.Background(), func(rm repo.RepositoryManager) error {
			return testErr
		})
		assert.Equal(t, testErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("transaction rollback error handling", func(t *testing.T) {
		testErr := errors.New("function error")
		rollbackErr := errors.New("rollback error")
		mock.ExpectBegin()
		mock.ExpectRollback().WillReturnError(rollbackErr)

		err := uow.Do(context.Background(), func(rm repo.RepositoryManager) error {
			return testErr
		})
		assert.Equal(t, rollbackErr, err) // Should return the rollback error
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("begin transaction error", func(t *testing.T) {
		beginErr := errors.New("begin error")
		mock.ExpectBegin().WillReturnError(beginErr)

		err := uow.Do(context.Background(), func(rm repo.RepositoryManager) error {
			return nil
		})
		assert.Equal(t, beginErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("repository manager provides correct repos", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		err := uow.Do(context.Background(), func(rm repo.RepositoryManager) error {
			assert.NotNil(t, rm.UserRepo())
			assert.NotNil(t, rm.AuthRepo())
			assert.NotNil(t, rm.ConnectionRepo())
			assert.NotNil(t, rm.MessageRepo())
			assert.NotNil(t, rm.NotificationRepo())
			assert.NotNil(t, rm.PasswordResetRepo())
			assert.NotNil(t, rm.PictureRepo())
			assert.NotNil(t, rm.RefreshTokenRepo())
			assert.NotNil(t, rm.UserTagRepo())
			assert.NotNil(t, rm.VerificationTokenRepo())
			assert.NotNil(t, rm.ProfileRepo())
			assert.NotNil(t, rm.ViewRepo())
			assert.NotNil(t, rm.LikeRepo())
			assert.NotNil(t, rm.BlockRepo())
			assert.NotNil(t, rm.UserDataRepo())
			return nil
		})
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
