package out

import (
	"context"

	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/user/domain/valueobject"
)

// UserAuthenticator is a port for authentication-related user operations
// This allows Auth domain to verify credentials without depending on User domain details
type UserAuthenticator interface {
	// VerifyCredentials checks if email+password combination is valid
	// Returns userID if valid, error otherwise
	VerifyCredentials(ctx context.Context, email valueobject.Email, plainPassword string) (userID string, err error)
}
