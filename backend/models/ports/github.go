package ports

import "context"

type GitHub interface {
	GetUserIDByUserName(ctx context.Context, name string) (string, error)
}
