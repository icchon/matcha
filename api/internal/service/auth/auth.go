package auth

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/icchon/matcha/api/internal/apperrors"
	"github.com/icchon/matcha/api/internal/domain/client"
	"github.com/icchon/matcha/api/internal/domain/entity"
	"github.com/icchon/matcha/api/internal/domain/repo"
	"github.com/icchon/matcha/api/internal/domain/service"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	uow               repo.UnitOfWork
	authRepo          repo.AuthQueryRepository
	userRepo          repo.UserQueryRepository
	refreshTokenRepo  repo.RefreshTokenQueryRepository
	passwordResetRepo repo.PasswordResetQueryRepository
	verificationRepo  repo.VerificationTokenQueryRepository
	hmacSecretKey     string
	jwtSigningKey     string
	mailService       service.MailService
	googleClient      client.OAuthClient
	githubClient      client.OAuthClient
}

var _ service.AuthService = (*authService)(nil)

func NewAuthService(
	uow repo.UnitOfWork,
	authRepo repo.AuthQueryRepository,
	userRepo repo.UserQueryRepository,
	refreshTokenRepo repo.RefreshTokenQueryRepository,
	passwordResetRepo repo.PasswordResetQueryRepository,
	verificationRepo repo.VerificationTokenQueryRepository,
	googleClient client.OAuthClient,
	githubClient client.OAuthClient,
	mailService service.MailService,
	hmacSecretKey string,
	jwtSigningKey string,
) *authService {
	return &authService{
		uow:               uow,
		authRepo:          authRepo,
		userRepo:          userRepo,
		refreshTokenRepo:  refreshTokenRepo,
		hmacSecretKey:     hmacSecretKey,
		jwtSigningKey:     jwtSigningKey,
		mailService:       mailService,
		verificationRepo:  verificationRepo,
		passwordResetRepo: passwordResetRepo,
		googleClient:      googleClient,
		githubClient:      githubClient,
	}
}

func (s *authService) LoginOAuth(ctx context.Context, code string, codeVerifier string, provider entity.AuthProvider) (a *entity.Auth, access string, refresh string, e error) {
	var oauthInfo *client.OAuthInfo
	var err error
	switch provider {
	case entity.ProviderGoogle:
		oauthInfo, err = s.googleClient.ExchangeCode(ctx, code, codeVerifier)
	case entity.ProviderGithub:
		oauthInfo, err = s.githubClient.ExchangeCode(ctx, code, codeVerifier)
	case entity.ProviderApple:
		return nil, "", "", apperrors.ErrNotImplemented
	case entity.ProviderFacebook:
		return nil, "", "", apperrors.ErrNotImplemented
	default:
		return nil, "", "", apperrors.ErrInternalServer
	}
	if err != nil {
		return nil, "", "", err
	}

	authes, err := s.authRepo.Query(ctx, &repo.AuthQuery{Provider: &provider, ProviderUID: &sql.NullString{String: oauthInfo.Sub, Valid: true}})
	if err != nil {
		return nil, "", "", apperrors.ErrInternalServer
	}
	if len(authes) > 1 {
		return nil, "", "", apperrors.ErrInternalServer
	}

	var auth *entity.Auth
	if len(authes) == 0 {
		err := s.uow.Do(ctx, func(m repo.RepositoryManager) error {
			user := &entity.User{}
			if err := m.UserRepo().Create(ctx, user); err != nil {
				return err
			}
			auth = &entity.Auth{
				UserID:      user.ID,
				Email:       sql.NullString{String: oauthInfo.Email, Valid: oauthInfo.EmailVerified},
				Provider:    provider,
				ProviderUID: sql.NullString{String: oauthInfo.Sub, Valid: true},
				IsVerified:  true,
			}
			return m.AuthRepo().Create(ctx, auth)
		})
		if err != nil {
			return nil, "", "", err
		}
	} else {
		auth = authes[0]
	}
	ac, re, err := s.IssueTokens(ctx, auth.UserID)
	return auth, ac, re, err
}

func (s *authService) ConfirmPassword(ctx context.Context, token string, password string) error {
	passwordToken, err := s.passwordResetRepo.Find(ctx, token)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	if passwordToken == nil {
		return apperrors.ErrUnauthorized
	}
	if passwordToken.ExpiresAt.Before(time.Now()) {
		return apperrors.ErrUnauthorized
	}
	passwordHash, err := HashPassword(password)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	auth, err := s.authRepo.Find(ctx, passwordToken.UserID, entity.ProviderLocal)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	if auth == nil {
		return apperrors.ErrNotFound
	}
	auth.PasswordHash = sql.NullString{String: passwordHash, Valid: true}
	return s.uow.Do(ctx, func(m repo.RepositoryManager) error {
		if err := m.AuthRepo().Update(ctx, auth); err != nil {
			return err
		}
		if err := m.PasswordResetRepo().Delete(ctx, token); err != nil {
			return err
		}
		return nil
	})
}

func (s *authService) SendPasswordResetEmail(ctx context.Context, email string) error {
	provider := entity.ProviderLocal
	auth, err := s.authRepo.Query(ctx, &repo.AuthQuery{Email: &sql.NullString{String: email, Valid: true}, Provider: &provider})
	if err != nil {
		return apperrors.ErrInternalServer
	}
	if len(auth) != 1 {
		return apperrors.ErrNotFound
	}
	token := GenerateEmailToken()
	if err := s.uow.Do(ctx, func(m repo.RepositoryManager) error {
		passwordReset := &entity.PasswordReset{
			UserID:    auth[0].UserID,
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour),
		}
		return m.PasswordResetRepo().Create(ctx, passwordReset)
	}); err != nil {
		return err
	}
	return s.mailService.SendPasswordResetEmail(ctx, email, token)
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	verification, err := s.verificationRepo.Find(ctx, token)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	if verification == nil {
		return apperrors.ErrUnauthorized
	}
	if verification.ExpiresAt.Before(time.Now()) {
		return apperrors.ErrUnauthorized
	}

	auth, err := s.authRepo.Find(ctx, verification.UserID, entity.ProviderLocal)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	if auth == nil {
		return apperrors.ErrNotFound
	}
	auth.IsVerified = true
	return s.uow.Do(ctx, func(m repo.RepositoryManager) error {
		if err := m.VerificationTokenRepo().Delete(ctx, token); err != nil {
			return err
		}
		if err := m.AuthRepo().Update(ctx, auth); err != nil {
			return err
		}
		return nil
	})
}

func (s *authService) SendVerificationEmail(ctx context.Context, email string, userID uuid.UUID) error {
	emailToken, err := s.IssueEMailToken(ctx, userID)
	if err != nil {
		return err
	}
	return s.mailService.SendVerificationEmail(ctx, email, "", emailToken)
}

func (s *authService) Signup(ctx context.Context, email string, password string) error {
	if !IsValidEmailFormat(email) {
		return apperrors.ErrInvalidInput
	}
	if !IsValidPasswordFormat(password) {
		return apperrors.ErrInvalidInput
	}
	provider := entity.ProviderLocal
	tokens, err := s.authRepo.Query(ctx, &repo.AuthQuery{Email: &sql.NullString{String: email, Valid: true}, Provider: &provider})
	if err != nil {
		log.Printf("query error: %v", err)
		return apperrors.ErrInternalServer
	}
	if len(tokens) > 0 {
		return apperrors.ErrInvalidInput
	}

	var id uuid.UUID
	if err := s.uow.Do(ctx, func(m repo.RepositoryManager) error {
		user := &entity.User{}
		log.Printf("Creating user for email: %s", email)
		if err := m.UserRepo().Create(ctx, user); err != nil {
			return err
		}
		log.Printf("User created with ID: %s", user.ID)
		id = user.ID
		passwordHash, err := HashPassword(password)
		if err != nil {
			log.Printf("password hash error: %v", err)
			return apperrors.ErrInternalServer
		}
		log.Printf("Password hashed for user ID: %s", user.ID)
		auth := &entity.Auth{
			UserID:       id,
			Email:        sql.NullString{String: email, Valid: true},
			Provider:     entity.ProviderLocal,
			PasswordHash: sql.NullString{String: passwordHash, Valid: true},
			IsVerified:   false,
		}
		log.Printf("Creating auth record for user ID: %s", user.ID)
		return m.AuthRepo().Create(ctx, auth)
	}); err != nil {
		return err
	}
	return s.SendVerificationEmail(ctx, email, id)
}

func (s *authService) IssueEMailToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token := GenerateEmailToken()
	tokens, err := s.verificationRepo.Query(ctx, &repo.VerificationTokenQuery{UserID: &userID})
	if err != nil {
		return "", apperrors.ErrInternalServer
	}
	if len(tokens) > 1 {
		return "", apperrors.ErrUnhandled
	} else if len(tokens) == 1 {
		tokens[0].ExpiresAt = time.Now().Add(time.Hour)
		tokens[0].Token = token
		if err := s.uow.Do(ctx, func(m repo.RepositoryManager) error {
			m.VerificationTokenRepo().Update(ctx, tokens[0])
			return nil
		}); err != nil {
			return "", err
		}
	} else if len(tokens) == 0 {
		if err := s.uow.Do(ctx, func(m repo.RepositoryManager) error {
			verificationToken := &entity.VerificationToken{
				UserID:    userID,
				Token:     token,
				ExpiresAt: time.Now().Add(time.Hour),
			}
			return m.VerificationTokenRepo().Create(ctx, verificationToken)
		}); err != nil {
			return "", err
		}
	}
	return token, nil
}

func (s *authService) IssueTokens(ctx context.Context, userID uuid.UUID) (string, string, error) {
	refreshToken, err := s.IssueRefreshToken(ctx, userID)
	if err != nil {
		log.Printf("issue refresh token error: %v", err)
		return "", "", err
	}
	accesToken, err := s.IssueAccessToken(ctx, refreshToken)
	if err != nil {
		log.Printf("issue access token error: %v", err)
		return "", "", err
	}
	return accesToken, refreshToken, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (a *entity.Auth, access string, refresh string, err error) {
	provider := entity.ProviderLocal
	auth, err := s.authRepo.Query(ctx, &repo.AuthQuery{Email: &sql.NullString{String: email, Valid: true}, Provider: &provider})
	if err != nil {
		log.Printf("query error: %v", err)
		return nil, "", "", apperrors.ErrInternalServer
	}
	if len(auth) != 1 {
		return nil, "", "", apperrors.ErrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(auth[0].PasswordHash.String), []byte(password)); err != nil {
		return nil, "", "", apperrors.ErrUnauthorized
	}

	accessToken, refreshToken, err := s.IssueTokens(ctx, auth[0].UserID)
	return auth[0], accessToken, refreshToken, err
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	refreshToken, err := s.refreshTokenRepo.Query(ctx, &repo.RefreshTokenQuery{UserID: &userID})
	if err != nil {
		return apperrors.ErrInternalServer
	}
	if len(refreshToken) != 1 {
		return apperrors.ErrNotFound
	}
	refreshToken[0].Revoked = true
	refreshToken[0].ExpiresAt = time.Now()

	user, err := s.userRepo.Find(ctx, refreshToken[0].UserID)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	if user == nil {
		return apperrors.ErrNotFound
	}
	user.LastConnection = sql.NullTime{
		Time: time.Now(),
	}

	return s.uow.Do(ctx, func(m repo.RepositoryManager) error {
		if err := m.RefreshTokenRepo().Update(ctx, refreshToken[0]); err != nil {
			return err
		}
		if err := m.UserRepo().Update(ctx, user); err != nil {
			return err
		}
		return nil
	})
}

func (s *authService) IssueAccessToken(ctx context.Context, oldRefreshToken string) (string, error) {
	token, err := s.VerifyRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		log.Printf("verify refresh token error: %v", err)
		return "", err
	}
	auths, err := s.authRepo.Query(ctx, &repo.AuthQuery{UserID: &token.UserID})
	if err != nil {
		return "", apperrors.ErrInternalServer
	}
	if len(auths) == 0 {
		return "", apperrors.ErrNotFound
	}
	auth := auths[0]
	accessToken, err := GenerateAccessToken(token.UserID, auth.IsVerified, auth.Provider, s.jwtSigningKey)
	if err != nil {
		log.Printf("generate access token error: %v", err)
		return "", apperrors.ErrInternalServer
	}
	return accessToken, nil
}

func (s *authService) IssueRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	refreshToken := GenerateRefreshToken()
	refreshTokenHash := HashTokenWithHMAC(refreshToken, s.hmacSecretKey)
	if err := s.uow.Do(ctx, func(m repo.RepositoryManager) error {
		refreshToken := &entity.RefreshToken{
			UserID:    userID,
			TokenHash: refreshTokenHash,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		return m.RefreshTokenRepo().Create(ctx, refreshToken)
	}); err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (s *authService) RefreshAccessToken(ctx context.Context, refreshTokenString string) (string, error) {
	refreshToken, err := s.VerifyRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return "", err
	}

	auths, err := s.authRepo.Query(ctx, &repo.AuthQuery{UserID: &refreshToken.UserID})
	if err != nil {
		return "", apperrors.ErrInternalServer
	}
	if len(auths) == 0 {
		return "", apperrors.ErrNotFound
	}
	auth := auths[0]

	newAccessToken, err := GenerateAccessToken(refreshToken.UserID, auth.IsVerified, auth.Provider, s.jwtSigningKey)
	if err != nil {
		return "", apperrors.ErrInternalServer
	}

	return newAccessToken, nil
}
