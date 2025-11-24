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

type userDataRepository struct {
	db DBTX
}

func NewUserDataRepository(db DBTX) repo.UserDataRepository {
	return &userDataRepository{db: db}
}

func (r *userDataRepository) Create(ctx context.Context, userData *entity.UserData) error {
	query := `
		INSERT INTO user_data (user_id, latitude, longitude, internal_score)
		VALUES (:user_id, :latitude, :longitude, :internal_score)
		RETURNING *
	`
	return r.db.QueryRowxContext(ctx, query, userData).StructScan(userData)
}

func (r *userDataRepository) Update(ctx context.Context, userData *entity.UserData) error {
	query := `
		UPDATE user_data SET
			latitude = :latitude,
			longitude = :longitude,
			internal_score = :internal_score
		WHERE user_id = :user_id
	`
	_, err := r.db.NamedExecContext(ctx, query, userData)
	return err
}

func (r *userDataRepository) Find(ctx context.Context, userID uuid.UUID) (*entity.UserData, error) {
	var userData entity.UserData
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

func (r *userDataRepository) Query(ctx context.Context, q *repo.UserDataQuery) ([]*entity.UserData, error) {
	query := "SELECT * FROM user_data WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *q.UserID)
		argCount++
	}
	if q.Latitude != nil {
		query += fmt.Sprintf(" AND latitude = $%d", argCount)
		args = append(args, *q.Latitude)
		argCount++
	}
	if q.Longitude != nil {
		query += fmt.Sprintf(" AND longitude = $%d", argCount)
		args = append(args, *q.Longitude)
		argCount++
	}
	if q.InternalScore != nil {
		query += fmt.Sprintf(" AND internal_score = $%d", argCount)
		args = append(args, *q.InternalScore)
		argCount++
	}

	var userData []*entity.UserData
	if err := r.db.SelectContext(ctx, &userData, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return userData, nil
}

func (r *userDataRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := "DELETE FROM user_data WHERE user_id = $1"
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
