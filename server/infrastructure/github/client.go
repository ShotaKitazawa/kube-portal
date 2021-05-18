package github

import (
	"context"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

type GitHubClient struct {
	logger *logrus.Logger
}

func NewGitHubClient(l *logrus.Logger) *GitHubClient {
	return &GitHubClient{l}
}

func (c *GitHubClient) GetUserIDByUserName(ctx context.Context, name string) (string, error) {
	client := github.NewClient(nil)
	user, _, err := client.Users.Get(ctx, name)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(*user.ID, 10), nil
}
