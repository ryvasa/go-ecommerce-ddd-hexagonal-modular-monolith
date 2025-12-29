package out

import (
	"context"

	userOut "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/out"
)

type UserReaderImpl struct {
	repo userOut.UserRepository
}

func NewUserReaderImpl(repo userOut.UserRepository) *UserReaderImpl {
	return &UserReaderImpl{
		repo: repo,
	}
}

func (r *UserReaderImpl) Exists(ctx context.Context, userID string) (bool, error) {
	return r.repo.ExistsByID(ctx, userID)
}
