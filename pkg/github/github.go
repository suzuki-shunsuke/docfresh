package github

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/go-github/v89/github"
	"github.com/suzuki-shunsuke/ghtkn-go-sdk/ghtkn"
	"golang.org/x/oauth2"
)

type Client struct {
	repo RepositoriesService
}

type RepositoriesService interface {
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error)
}

type (
	RepositoryContentGetOptions = github.RepositoryContentGetOptions
)

func New(ctx context.Context, logger *slog.Logger, token string) (*Client, error) {
	hc, err := getHTTPClient(ctx, logger, token)
	if err != nil {
		return nil, fmt.Errorf("get HTTP client: %w", err)
	}
	gh, err := github.NewClient(github.WithHTTPClient(hc))
	if err != nil {
		return nil, fmt.Errorf("create a GitHub client: %w", err)
	}
	return &Client{
		repo: gh.Repositories,
	}, nil
}

func getHTTPClient(ctx context.Context, logger *slog.Logger, token string) (*http.Client, error) {
	ts, err := getTokenSource(logger, token)
	if err != nil {
		return nil, fmt.Errorf("get token source: %w", err)
	}
	if ts == nil {
		return http.DefaultClient, nil
	}
	return oauth2.NewClient(ctx, ts), nil
}

func getTokenSource(logger *slog.Logger, token string) (oauth2.TokenSource, error) {
	if token != "" {
		return oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		), nil
	}
	ghtknEnabled, err := ghtkn.Enabled(&ghtkn.InputEnabled{
		Envs: []string{"DOCFRESH_GHTKN_ENABLED"},
	})
	if err != nil {
		return nil, fmt.Errorf("check if ghtkn is enabled: %w", err)
	}
	if !ghtknEnabled {
		return nil, nil //nolint:nilnil
	}
	client, err := ghtkn.New()
	if err != nil {
		return nil, fmt.Errorf("create a ghtkn client: %w", err)
	}
	return client.TokenSource(logger, &ghtkn.InputGet{}), nil
}

func GetGitHubTokenFromEnv() string {
	for _, key := range []string{"DOCFRESH_GITHUB_TOKEN", "GITHUB_TOKEN"} {
		s := os.Getenv(key)
		if s != "" {
			return s
		}
	}
	return ""
}
