package in

import "context"

// LoginCommand represents login credentials
type LoginCommand struct {
	Email    string
	Password string
}

// AuthUsecase defines authentication use cases
type AuthUsecase interface {
	Login(ctx context.Context, cmd LoginCommand) (accessToken string, err error)
	// Future: Logout, RefreshToken, etc.
}
