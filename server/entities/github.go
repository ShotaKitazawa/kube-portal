package entities

import "context"

type GitHubPort interface {
	GetUserIDByUserName(ctx context.Context, name string) (string, error)
}
