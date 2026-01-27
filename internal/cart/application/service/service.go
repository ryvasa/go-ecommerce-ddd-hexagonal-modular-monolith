package service

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/application/port/out"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/cart/domain/entity"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/apperror"
	sharedEvent "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/event"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/transaction"
)

type CartService struct {
	repo       out.CartRepository
	userReader out.UserReader
	txManager  transaction.Manager
}

func NewCartService(
	repo out.CartRepository,
	userReader out.UserReader,
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
		// Check if user exists (cross-module check, use shared apperror)
		exists, err := s.userReader.Exists(txCtx, userID)
		if err != nil {
			return err
		}
		if !exists {
			return apperror.NotFound("user not found")
		}

		// Create cart using domain factory - validation happens in entity
		cart, err := entity.NewCart(userID, title)
		if err != nil {
			return err
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
