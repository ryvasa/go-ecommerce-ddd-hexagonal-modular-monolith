package out

import "context"

type UserReader interface {
	Exists(ctx context.Context, userID string) (bool, error)
}
