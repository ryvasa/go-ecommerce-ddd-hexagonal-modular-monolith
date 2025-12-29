package casbin

import (
	lib "github.com/casbin/casbin/v2"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/domain/authorization"
	"gorm.io/gorm"
)

// ProvideEnforcer creates a Casbin enforcer for dependency injection.
// It will use database adapter if db is provided, otherwise uses file-based policies.
func ProvideEnforcer(db *gorm.DB) (*lib.Enforcer, error) {
	// For now, we'll pass nil to use file-based policies
	// In production, you can enable database adapter by passing db
	return NewEnforcer(nil)
}

// ProvideAuthorizer creates an Authorizer implementation for dependency injection.
// This binds the infrastructure adapter to the domain interface.
func ProvideAuthorizer(enforcer *lib.Enforcer) authorization.Authorizer {
	return NewAdapter(enforcer)
}
