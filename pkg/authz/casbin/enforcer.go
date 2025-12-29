package casbin

import lib "github.com/casbin/casbin/v2"

type Enforcer struct {
	enforcer *lib.Enforcer
}

func New(enforcer *lib.Enforcer) *Enforcer {
	return &Enforcer{enforcer: enforcer}
}

func (e *Enforcer) Enforce(sub, obj, act string) (bool, error) {
	return e.enforcer.Enforce(sub, obj, act)
}
