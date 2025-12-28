package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/domain/entity"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/apperror"
	sharedEvent "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/event"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/transaction"
	userIn "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
)

type CartService struct {
	repo       out.CartRepository
	userReader userIn.UserReader
	txManager  transaction.Manager
}

func NewCartService(
	repo out.CartRepository,
	userReader userIn.UserReader,
	txManager transaction.Manager,
	eventPublisher sharedEvent.Publisher,
) in.CartUsecase {
	return &CartService{
		repo:       repo,
		userReader: userReader,
		txManager:  txManager,
	}
}

func (s *CartService) Create(ctx context.Context, userID, title string) error {
	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		exists, err := s.userReader.Exists(txCtx, userID)
		if err != nil {
			return err
		}
		if !exists {
			return apperror.NotFound("user not found")
		}

		cart := &entity.Cart{
			ID:     uuid.NewString(),
			UserID: userID,
			Title:  title,
			Done:   false,
		}

		err = s.repo.Save(txCtx, cart)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *CartService) List(ctx context.Context) ([]*entity.Cart, error) {
	return s.repo.FindAll(ctx)
}
