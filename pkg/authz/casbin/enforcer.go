package casbin

import (
	lib "github.com/casbin/casbin/v2"
	"github.com/ryvasa/go-ddd-hexagonal-modular-monolith/pkg/config"
)

type Enforcer struct {
	enforcer *lib.Enforcer
}

func NewCasbinEnforcer(cfg config.CasbinConfig) (*Enforcer, error) {
	e, err := lib.NewEnforcer(cfg.ModelPath, cfg.PolicyPath)
	if err != nil {
		return nil, err
	}

	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}

	return &Enforcer{enforcer: e}, nil
}

func (e *Enforcer) Enforce(sub, obj, act string) (bool, error) {
	return e.enforcer.Enforce(sub, obj, act)
}
