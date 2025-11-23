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

type userProfileRepository struct {
	db DBTX
}

func NewUserProfileRepository(db DBTX) repo.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Create(ctx context.Context, userProfile *entity.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, first_name, last_name, username, gender, sexual_preference, biography, location_name)
		VALUES (:user_id, :first_name, :last_name, :username, :gender, :sexual_preference, :biography, :location_name)
		RETURNING *
	`
	return r.db.QueryRowxContext(ctx, query, userProfile).StructScan(userProfile)
}

func (r *userProfileRepository) Update(ctx context.Context, userProfile *entity.UserProfile) error {
	query := `
		UPDATE user_profiles SET
			first_name = :first_name,
			last_name = :last_name,
			username = :username,
			gender = :gender,
			sexual_preference = :sexual_preference,
			biography = :biography,
			fame_rating = :fame_rating,
			location_name = :location_name
		WHERE user_id = :user_id
		RETURNING *
	`
	_, err := r.db.NamedExecContext(ctx, query, userProfile)
	return err
}

func (r *userProfileRepository) Find(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	var userProfile entity.UserProfile
	query := "SELECT * FROM user_profiles WHERE user_id = $1"
	err := r.db.GetContext(ctx, &userProfile, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &userProfile, nil
}

func (r *userProfileRepository) Query(ctx context.Context, q *repo.UserProfileQuery) ([]*entity.UserProfile, error) {
	query := "SELECT * FROM user_profiles WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *q.UserID)
		argCount++
	}
	if q.FirstName != nil {
		query += fmt.Sprintf(" AND first_name = $%d", argCount)
		args = append(args, *q.FirstName)
		argCount++
	}
	if q.LastName != nil {
		query += fmt.Sprintf(" AND last_name = $%d", argCount)
		args = append(args, *q.LastName)
		argCount++
	}
	if q.Username != nil {
		query += fmt.Sprintf(" AND username = $%d", argCount)
		args = append(args, *q.Username)
		argCount++
	}
	if q.Gender != nil {
		query += fmt.Sprintf(" AND gender = $%d", argCount)
		args = append(args, *q.Gender)
		argCount++
	}
	if q.SexualPreference != nil {
		query += fmt.Sprintf(" AND sexual_preference = $%d", argCount)
		args = append(args, *q.SexualPreference)
		argCount++
	}
	if q.Biography != nil {
		query += fmt.Sprintf(" AND biography LIKE $%d", argCount)
		args = append(args, "%"+*q.Biography+"%")
		argCount++
	}
	if q.FameRating != nil {
		query += fmt.Sprintf(" AND fame_rating = $%d", argCount)
		args = append(args, *q.FameRating)
		argCount++
	}
	if q.LocationName != nil {
		query += fmt.Sprintf(" AND location_name = $%d", argCount)
		args = append(args, *q.LocationName)
		argCount++
	}

	var userProfiles []*entity.UserProfile
	if err := r.db.SelectContext(ctx, &userProfiles, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return userProfiles, nil
}

func (r *userProfileRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := "DELETE FROM user_profiles WHERE user_id = $1"
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
