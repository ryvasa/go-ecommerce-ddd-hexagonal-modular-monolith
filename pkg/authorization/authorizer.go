package authorization

// Authorizer defines the port (interface) for authorization operations.
// This is an infrastructure-level interface that is implementation-agnostic.
// Infrastructure adapters (like Casbin) will implement this interface.
type Authorizer interface {
	// Enforce checks if a subject can perform an action on an object
	Enforce(subject, object, action string) (bool, error)

	// AddPolicy adds a new authorization policy
	AddPolicy(subject, object, action string) error

	// RemovePolicy removes an authorization policy
	RemovePolicy(subject, object, action string) error

	// AddRoleForUser adds a role for a user
	AddRoleForUser(user, role string) error

	// RemoveRoleForUser removes a role from a user
	RemoveRoleForUser(user, role string) error
}
