package out

import (
	"context"

	authPort "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/auth/application/port/out"
	userPort "github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/application/port/in"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

// UserAuthenticatorImpl adapts User module's UserReader to Auth module's needs
type UserAuthenticatorImpl struct {
	userReader userPort.UserReader
}

func NewUserAuthenticatorImpl(ur userPort.UserReader) authPort.UserAuthenticator {
	return &UserAuthenticatorImpl{userReader: ur}
}

func (a *UserAuthenticatorImpl) VerifyCredentials(
	ctx context.Context,
	email string,
	plainPassword string,
) (string, error) {
	// Conversion from primitive to value object happens here in the adapter
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return "", err
	}
	return a.userReader.VerifyCredentials(ctx, emailVO, plainPassword)
}
