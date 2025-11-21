package oauth

import (
	"context"
	"github.com/icchon/matcha/api/internal/domain/repo"
)

type githubClient struct {
	githubClientID     string
	githubClientSecret string
	githubRedirectURL  string
}

var _ repo.OAuthClient = (*githubClient)(nil)

func NewGithubClient(githubClientID, githubClientSecret, githubRedirectURL string) *githubClient {
	return &githubClient{
		githubClientID:     githubClientID,
		githubClientSecret: githubClientSecret,
		githubRedirectURL:  githubRedirectURL,
	}
}

func (h *githubClient) ExchangeCode(ctx context.Context, code string, codeVerifer string) (*repo.OAuthInfo, error) {
	return nil, nil
}
