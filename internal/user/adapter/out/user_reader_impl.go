package out

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain"
)

type UserReaderImpl struct {
	repo domain.UserRepository
}

func NewUserReaderImpl(repo domain.UserRepository) *UserReaderImpl {
	return &UserReaderImpl{
		repo: repo,
	}
}

func (r *UserReaderImpl) Exists(ctx context.Context, userID string) (bool, error) {
	return r.repo.ExistsByID(ctx, userID)
}
