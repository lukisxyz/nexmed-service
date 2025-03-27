package config

import "fmt"

type redisConfig struct {
	Host string `yaml:"host" json:"host"`
	Port uint   `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
}

func (l redisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

func defaultRedisConfig() redisConfig {
	return redisConfig{
		Host: "127.0.0.1",
		Port: 6379,
		Password: "",
	}
}

func (l *redisConfig) loadFromEnv() {
	loadEnvStr("REDIS_HOST", &l.Host)
	loadEnvUint("REDIS_PORT", &l.Port)
	loadEnvStr("REDIS_PASSWORD", &l.Password)
}
