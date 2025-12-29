package config

type CasbinConfig struct {
	ModelPath  string
	PolicyPath string
}

func NewCasbinConfig() CasbinConfig {
	return CasbinConfig{
		ModelPath:  GetEnv("CASBIN_MODEL_PATH", true),
		PolicyPath: GetEnv("CASBIN_POLICY_PATH", true),
	}
}
