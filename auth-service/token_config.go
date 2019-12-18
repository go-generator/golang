package auth_service

type TokenConfig struct {
	Secret  string `mapstructure:"secret"`
	Expires uint64 `mapstructure:"expires"`
}
