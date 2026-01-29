package casbin

import (
	lib "github.com/casbin/casbin/v2"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/authorization"
	"gorm.io/gorm"
)

// ProvideEnforcer creates a Casbin enforcer for dependency injection.
// It uses database adapter for persistent policy storage.
func ProvideEnforcer(db *gorm.DB) (*lib.Enforcer, error) {
	return NewEnforcer(db)
}

// ProvideAuthorizer creates an Authorizer implementation for dependency injection.
// This binds the infrastructure adapter to the domain interface.
func ProvideAuthorizer(enforcer *lib.Enforcer) authorization.Authorizer {
	return NewAdapter(enforcer)
}
