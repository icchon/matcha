package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain"
)

type userDataRepository struct {
	db DBTX
}

func NewUserDataRepository(db DBTX) domain.UserDataRepository {
	return &userDataRepository{db: db}
}

func (r *userDataRepository) Save(ctx context.Context, userData *domain.UserData) error {
	query := `
		INSERT INTO user_data (user_id, latitude, longitude, internal_score)
		VALUES (:user_id, :latitude, :longitude, :internal_score)
		ON CONFLICT (user_id) DO UPDATE SET
			latitude = :latitude,
			longitude = :longitude,
			internal_score = :internal_score
	`
	_, err := r.db.NamedExecContext(ctx, query, userData)
	return err
}

func (r *userDataRepository) Find(ctx context.Context, userID uuid.UUID) (*domain.UserData, error) {
	var userData domain.UserData
	query := "SELECT * FROM user_data WHERE user_id = $1"
	err := r.db.GetContext(ctx, &userData, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &userData, nil
}

func (r *userDataRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := "DELETE FROM user_data WHERE user_id = $1"
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
