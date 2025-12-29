package casbin

import lib "github.com/casbin/casbin/v2"

func ProvideAuthorizer(e *lib.Enforcer) *Enforcer {
	return New(e)
}
