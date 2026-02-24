package profile

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProfileService_RecommendProfiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	profileRepo := mock.NewMockUserProfileRepository(ctrl)
	userDataRepo := mock.NewMockUserDataRepository(ctrl)
	userTagRepo := mock.NewMockUserTagRepository(ctrl)
	blockRepo := mock.NewMockBlockRepository(ctrl)
	// other repos and services can be mocked as needed

	profileSvc := NewProfileService(nil, profileRepo, nil, nil, nil, nil, blockRepo, nil, userTagRepo, userDataRepo)

	selfUserID := uuid.New()
	candidateUserID1 := uuid.New()
	candidateUserID2 := uuid.New()

	selfData := &entity.UserData{
		UserID:    selfUserID,
		Latitude:  sql.NullFloat64{Float64: 35.68, Valid: true},
		Longitude: sql.NullFloat64{Float64: 139.76, Valid: true},
	}

	selfTags := []*entity.Tag{
		{ID: 1, Name: "music"},
		{ID: 2, Name: "sports"},
	}

	candidate1 := &entity.UserProfile{
		UserID:   candidateUserID1,
		Distance: sql.NullFloat64{Float64: 3000, Valid: true}, // 3km
	}
	candidate1Tags := []*entity.Tag{
		{ID: 1, Name: "music"}, // Common tag
	}

	candidate2 := &entity.UserProfile{
		UserID:   candidateUserID2,
		Distance: sql.NullFloat64{Float64: 10000, Valid: true}, // 10km
	}
	candidate2Tags := []*entity.Tag{
		{ID: 2, Name: "sports"}, // Common tag
	}

	// Expectations
	userDataRepo.EXPECT().Find(gomock.Any(), selfUserID).Return(selfData, nil)
	profileRepo.EXPECT().Query(gomock.Any(), gomock.Any()).Return([]*entity.UserProfile{candidate1, candidate2}, nil)
	userTagRepo.EXPECT().Query(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, q *repo.UserTagQuery) ([]*entity.Tag, error) {
		if *q.UserID == selfUserID {
			return selfTags, nil
		}
		if *q.UserID == candidateUserID1 {
			return candidate1Tags, nil
		}
		if *q.UserID == candidateUserID2 {
			return candidate2Tags, nil
		}
		return nil, nil
	}).AnyTimes()

	recs, err := profileSvc.RecommendProfiles(context.Background(), selfUserID)

	assert.NoError(t, err)
	assert.Len(t, recs, 2)
	// candidate1 has 1 common tag (10) + distance < 5km (20) = 30
	// candidate2 has 1 common tag (10) + distance < 20km (10) = 20
	// So candidate1 should be first
	assert.Equal(t, candidateUserID1, recs[0].UserID)
	assert.Equal(t, candidateUserID2, recs[1].UserID)
}
