package casbin

import (
	lib "github.com/casbin/casbin/v2"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/internal/shared/domain/authorization"
)

// Adapter implements the domain authorization.Authorizer interface using Casbin.
// This is the infrastructure adapter that translates domain operations to Casbin calls.
type Adapter struct {
	enforcer *lib.Enforcer
}

// NewAdapter creates a new Casbin adapter
func NewAdapter(enforcer *lib.Enforcer) *Adapter {
	return &Adapter{enforcer: enforcer}
}

// Enforce checks if a subject can perform an action on an object
func (a *Adapter) Enforce(subject, object, action string) (bool, error) {
	return a.enforcer.Enforce(subject, object, action)
}

// AddPolicy adds a new authorization policy
func (a *Adapter) AddPolicy(subject, object, action string) error {
	_, err := a.enforcer.AddPolicy(subject, object, action)
	return err
}

// RemovePolicy removes an authorization policy
func (a *Adapter) RemovePolicy(subject, object, action string) error {
	_, err := a.enforcer.RemovePolicy(subject, object, action)
	return err
}

// AddRoleForUser adds a role for a user
func (a *Adapter) AddRoleForUser(user, role string) error {
	_, err := a.enforcer.AddRoleForUser(user, role)
	return err
}

// RemoveRoleForUser removes a role from a user
func (a *Adapter) RemoveRoleForUser(user, role string) error {
	_, err := a.enforcer.DeleteRoleForUser(user, role)
	return err
}

// Ensure Adapter implements authorization.Authorizer interface
var _ authorization.Authorizer = (*Adapter)(nil)
