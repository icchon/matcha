package oauth

import (
	"context"
	"github.com/icchon/matcha/api/internal/domain/client"
)

type githubClient struct {
	githubClientID     string
	githubClientSecret string
	githubRedirectURL  string
}

var _ client.OAuthClient = (*githubClient)(nil)

func NewGithubClient(githubClientID, githubClientSecret, githubRedirectURL string) *githubClient {
	return &githubClient{
		githubClientID:     githubClientID,
		githubClientSecret: githubClientSecret,
		githubRedirectURL:  githubRedirectURL,
	}
}

func (h *githubClient) ExchangeCode(ctx context.Context, code string, codeVerifer string) (*client.OAuthInfo, error) {
	return nil, nil
}
