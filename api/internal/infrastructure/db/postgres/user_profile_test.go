package postgres

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserProfileRepository_Query(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	r := NewUserProfileRepository(db)

	testCases := []struct {
		name          string
		query         *repo.UserProfileQuery
		expectedQuery string
		expectedArgs  []interface{}
		expectedErr   error
	}{
		{
			name: "Filter by AgeMin and AgeMax",
			query: &repo.UserProfileQuery{
				AgeMin: func(i int) *int { return &i }(20),
				AgeMax: func(i int) *int { return &i }(30),
			},
			expectedQuery: `SELECT \* FROM user_profiles WHERE 1=1 AND date_part\('year', age\(birthday\)\) >= \$1 AND date_part\('year', age\(birthday\)\) <= \$2`,
			expectedArgs:  []interface{}{20, 30},
		},
		{
			name: "Filter by FameMin and FameMax",
			query: &repo.UserProfileQuery{
				FameMin: func(i int32) *int32 { return &i }(100),
				FameMax: func(i int32) *int32 { return &i }(500),
			},
			expectedQuery: `SELECT \* FROM user_profiles WHERE 1=1 AND fame_rating >= \$1 AND fame_rating <= \$2`,
			expectedArgs:  []interface{}{int32(100), int32(500)},
		},
		{
			name: "Filter by Gender",
			query: &repo.UserProfileQuery{
				Gender: func(g entity.Gender) *entity.Gender { return &g }("female"),
			},
			expectedQuery: `SELECT \* FROM user_profiles WHERE 1=1 AND gender = \$1`,
			expectedArgs:  []interface{}{entity.Gender("female")},
		},
		{
			name: "Exclude UserID",
			query: &repo.UserProfileQuery{
				ExcludeUserID: &userID1,
			},
			expectedQuery: `SELECT \* FROM user_profiles WHERE 1=1 AND user_id != \$1`,
			expectedArgs:  []interface{}{userID1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"user_id", "first_name", "last_name", "username", "gender", "sexual_preference", "birthday", "occupation", "biography", "fame_rating", "location_name"}).
				AddRow(userID2, "Test", "User", "testuser", "male", "bisexual", time.Now().AddDate(-25, 0, 0), "developer", "bio", 200, "Test Location")

			driverArgs := make([]driver.Value, len(tc.expectedArgs))
			for i, v := range tc.expectedArgs {
				driverArgs[i] = v
			}

			mock.ExpectQuery(tc.expectedQuery).
				WithArgs(driverArgs...).
				WillReturnRows(rows)

			_, err := r.Query(context.Background(), tc.query)

			assert.Equal(t, tc.expectedErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
