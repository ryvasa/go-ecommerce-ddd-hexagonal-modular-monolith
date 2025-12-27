package in

import "context"

type UserUsecase interface {
	Create(ctx context.Context, email string) error
}
