package casbin

import (
	lib "github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
)

func NewEnforcer() (*lib.Enforcer, error) {
	modelBytes, err := FS.ReadFile("model.conf")
	if err != nil {
		return nil, err
	}

	m := casbinModel.NewModel()
	err = m.LoadModelFromText(string(modelBytes))
	if err != nil {
		return nil, err
	}

	e, err := lib.NewEnforcer(m)
	if err != nil {
		return nil, err
	}

	// kalau policy file
	_, _ = FS.ReadFile("policy.csv") // optional
	_ = e.LoadPolicy()

	return e, nil
}
