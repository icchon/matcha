package profile

import (
	"context"
	"database/sql"
	"log"
	"sort"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/client"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/domain/service"
)

type profileService struct {
	uow          repo.UnitOfWork
	profileRepo  repo.UserProfileRepository
	pictureRepo  repo.PictureQueryRepository
	viewRepo     repo.ViewQueryRepository
	likeRepo     repo.LikeQueryRepository
	blockRepo    repo.BlockQueryRepository
	notifSvc     service.NotificationService
	fileClient   client.FileClient
	userTagRepo  repo.UserTagRepository
	userDataRepo repo.UserDataRepository
}

var _ service.ProfileService = (*profileService)(nil)

func NewProfileService(uow repo.UnitOfWork, profileRepo repo.UserProfileRepository, fileClient client.FileClient, pictureRepo repo.PictureQueryRepository, viewRepo repo.ViewQueryRepository, likeRepo repo.LikeQueryRepository, blockRepo repo.BlockQueryRepository, notifSvc service.NotificationService, userTagRepo repo.UserTagRepository, userDataRepo repo.UserDataRepository) *profileService {
	return &profileService{uow: uow, profileRepo: profileRepo, fileClient: fileClient, pictureRepo: pictureRepo, viewRepo: viewRepo, likeRepo: likeRepo, blockRepo: blockRepo, notifSvc: notifSvc, userTagRepo: userTagRepo, userDataRepo: userDataRepo}
}

func (s *profileService) CreateProfile(ctx context.Context, profile *entity.UserProfile) (*entity.UserProfile, error) {
	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		return rm.ProfileRepo().Create(ctx, profile)
	}); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	return profile, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, userID uuid.UUID, profile *entity.UserProfile) (*entity.UserProfile, error) {
	target, err := s.profileRepo.Find(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	if target == nil {
		return nil, apperrors.ErrNotFound
	}

	{
		if profile.FirstName.Valid {
			target.FirstName = profile.FirstName
		}
		if profile.LastName.Valid {
			target.LastName = profile.LastName
		}
		if profile.Username.Valid {
			target.Username = profile.Username
		}
		if profile.Gender.Valid {
			target.Gender = profile.Gender
		}
		if profile.SexualPreference.Valid {
			target.SexualPreference = profile.SexualPreference
		}
		if profile.Birthday.Valid {
			target.Birthday = profile.Birthday
		}
		if profile.Occupation.Valid {
			target.Occupation = profile.Occupation
		}
		if profile.Biography.Valid {
			target.Biography = profile.Biography
		}
		if profile.LocationName.Valid {
			target.LocationName = profile.LocationName
		}
	}

	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		return rm.ProfileRepo().Update(ctx, target)
	}); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	return target, nil
}

func (s *profileService) FindProfile(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	profile, err := s.profileRepo.Find(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, apperrors.ErrNotFound
	}
	// Calculate fame rating
	fameRating, err := s.calculateFameRating(ctx, userID)
	if err != nil {
		return nil, err
	}
	profile.FameRating = sql.NullInt32{Int32: fameRating, Valid: true}

	return profile, nil
}

func (s *profileService) ViewProfile(ctx context.Context, viewerID, viewedID uuid.UUID) error {
	if err := s.uow.Do(ctx, func(rm repo.RepositoryManager) error {
		view := &entity.View{
			ViewerID: viewerID,
			ViewedID: viewedID,
		}
		return rm.ViewRepo().Create(ctx, view)
	}); err != nil {
		return err
	}
	if _, err := s.notifSvc.CreateAndSendNotification(ctx, viewerID, viewedID, entity.NotifView); err != nil {
		return err
	}
	return nil
}

func (s *profileService) FindWhoViewedMeList(ctx context.Context, userID uuid.UUID) ([]*entity.View, error) {
	views, err := s.viewRepo.Query(ctx, &repo.ViewQuery{ViewedID: &userID})
	if err != nil {
		return nil, err
	}
	return views, nil
}

func (s *profileService) FindWhoLikedMeList(ctx context.Context, userID uuid.UUID) ([]*entity.Like, error) {
	likes, err := s.likeRepo.Query(ctx, &repo.LikeQuery{LikedID: &userID})
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (s *profileService) ListProfiles(ctx context.Context, selfUserID uuid.UUID, lat, lon, dist *float64, ageMin, ageMax *int, gender *entity.Gender) ([]*entity.UserProfile, error) {
	profiles, err := s.profileRepo.Query(ctx, &repo.UserProfileQuery{
		ExcludeUserID: &selfUserID,
		Latitude:      lat,
		Longitude:     lon,
		Distance:      dist,
		AgeMin:        ageMin,
		AgeMax:        ageMax,
		Gender:        gender,
	})
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	for _, profile := range profiles {
		fameRating, err := s.calculateFameRating(ctx, profile.UserID)
		if err != nil {
			// Log the error but continue processing other profiles
			log.Printf("Failed to calculate fame rating for user %s: %v", profile.UserID, err)
			profile.FameRating = sql.NullInt32{Int32: 0, Valid: true} // Default or error value
		} else {
			profile.FameRating = sql.NullInt32{Int32: fameRating, Valid: true}
		}
	}
	return profiles, nil
}

func (s *profileService) RecommendProfiles(ctx context.Context, selfUserID uuid.UUID) ([]*entity.UserProfile, error) {
	// 1. Get current user's data for scoring
	selfData, err := s.userDataRepo.Find(ctx, selfUserID)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	if selfData == nil {
		return []*entity.UserProfile{}, nil // No location data, no recommendations
	}

	// 2. DBから位置情報で絞り込んだ候補者リストを取得 (e.g., 50km radius)
	dist := 50.0
	candidateProfiles, err := s.profileRepo.Query(ctx, &repo.UserProfileQuery{
		ExcludeUserID: &selfUserID,
		Latitude:      &selfData.Latitude.Float64,
		Longitude:     &selfData.Longitude.Float64,
		Distance:      &dist,
	})
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	if len(candidateProfiles) == 0 {
		return []*entity.UserProfile{}, nil
	}

	// 3. 現在のユーザーのタグを取得
	selfTags, err := s.userTagRepo.Query(ctx, &repo.UserTagQuery{UserID: &selfUserID})
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	selfTagMap := make(map[int32]struct{})
	for _, tag := range selfTags {
		selfTagMap[tag.ID] = struct{}{}
	}

	// 4. スコアリングとソートを実行
	type scoredProfile struct {
		Profile *entity.UserProfile
		Score   int
	}
	scoredProfiles := make([]scoredProfile, 0, len(candidateProfiles))

	for _, profile := range candidateProfiles {
		candidateTags, err := s.userTagRepo.Query(ctx, &repo.UserTagQuery{UserID: &profile.UserID})
		if err != nil {
			candidateTags = []*entity.Tag{} // Treat as no tags for scoring
		}

		score := s.calculateScore(selfTagMap, candidateTags, profile.Distance)
		scoredProfiles = append(scoredProfiles, scoredProfile{Profile: profile, Score: score})
	}

	// スコアに基づいて降順にソート
	sort.Slice(scoredProfiles, func(i, j int) bool {
		return scoredProfiles[i].Score > scoredProfiles[j].Score
	})

	// ソートされたプロファイルのスライスを構築
	sortedProfiles := make([]*entity.UserProfile, 0, len(scoredProfiles))
	for _, sp := range scoredProfiles {
		sortedProfiles = append(sortedProfiles, sp.Profile)
	}

	return sortedProfiles, nil
}

// calculateScore は、共通タグと距離に基づいてユーザーのスコアを計算します
func (s *profileService) calculateScore(selfTagMap map[int32]struct{}, candidateTags []*entity.Tag, distance sql.NullFloat64) int {
	score := 0

	// 共通タグのスコア
	for _, tag := range candidateTags {
		if _, ok := selfTagMap[tag.ID]; ok {
			score += 10 // 共通タグ1つにつき10点
		}
	}

	// 距離のスコア
	if distance.Valid {
		distanceInKm := distance.Float64 / 1000 // メートルからキロメートルに変換
		if distanceInKm < 5 {
			score += 20 // 5km以内
		} else if distanceInKm < 20 {
			score += 10 // 20km以内
		} else if distanceInKm < 50 {
			score += 5 // 50km以内
		}
	}

	return score
}

// calculateFameRating は、ユーザーのFame Ratingを計算します。
// いいね数 + 閲覧数 - ブロック数の重み付けスコア
func (s *profileService) calculateFameRating(ctx context.Context, userID uuid.UUID) (int32, error) {
	likes, err := s.likeRepo.Query(ctx, &repo.LikeQuery{LikedID: &userID})
	if err != nil {
		return 0, err
	}
	views, err := s.viewRepo.Query(ctx, &repo.ViewQuery{ViewedID: &userID})
	if err != nil {
		return 0, err
	}
	blocks, err := s.blockRepo.Query(ctx, &repo.BlockQuery{BlockedID: &userID})
	if err != nil {
		return 0, err
	}

	// Simple weighting for now
	fameRating := (len(likes) * 5) + (len(views) * 1) - (len(blocks) * 10)
	if fameRating < 0 {
		fameRating = 0
	}

	return int32(fameRating), nil
}
