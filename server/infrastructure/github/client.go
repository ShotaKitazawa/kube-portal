package github

import (
	"context"
	"strconv"

	"github.com/ShotaKitazawa/kube-portal/server/models/ports"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

type Client struct {
	logger *logrus.Logger
}

var _ ports.GitHub = (*Client)(nil)

func NewGitHubClient(l *logrus.Logger) *Client {
	return &Client{l}
}

func (c *Client) GetUserIDByUserName(ctx context.Context, name string) (string, error) {
	client := github.NewClient(nil)
	user, _, err := client.Users.Get(ctx, name)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(*user.ID, 10), nil
}
