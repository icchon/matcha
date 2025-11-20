package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type userProfileRepository struct {
	db DBTX
}

func NewUserProfileRepository(db DBTX) domain.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Save(ctx context.Context, userProfile *domain.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, first_name, last_name, username, gender, sexual_preference, biography, fame_rating, location_name)
		VALUES (:user_id, :first_name, :last_name, :username, :gender, :sexual_preference, :biography, :fame_rating, :location_name)
		ON CONFLICT (user_id) DO UPDATE SET
			first_name = :first_name,
			last_name = :last_name,
			username = :username,
			gender = :gender,
			sexual_preference = :sexual_preference,
			biography = :biography,
			fame_rating = :fame_rating,
			location_name = :location_name
	`
	_, err := r.db.NamedExecContext(ctx, query, userProfile)
	return err
}

func (r *userProfileRepository) Find(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error) {
	var userProfile domain.UserProfile
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

func (r *userProfileRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := "DELETE FROM user_profiles WHERE user_id = $1"
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
