package config

type jwtConfig struct {
	Secret string `yaml:"secret" json:"secret"`
}

func defaultJwtConfig() jwtConfig {
	return jwtConfig{
		Secret: "CU6Y/cpppajGGIEDrMNkHmqvzA+5ShS5XStYqGzv++rnIhjMo6PfC/Aez2zlbWT3",
	}
}

func (l *jwtConfig) loadFromEnv() {
	loadEnvStr("JWT_SECRET", &l.Secret)
}
