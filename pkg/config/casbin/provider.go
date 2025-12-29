package casbin

import (
	lib "github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
)

func ProvideCasbinEnforcer() (*lib.Enforcer, error) {
	modelBytes, err := FS.ReadFile("model.conf")
	if err != nil {
		return nil, err
	}

	m := casbinModel.NewModel()
	if err := m.LoadModelFromText(string(modelBytes)); err != nil {
		return nil, err
	}

	return lib.NewEnforcer(m)
}
