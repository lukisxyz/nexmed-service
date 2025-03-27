package config

import (
	"io"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

func loadEnvStr(key string, result *string) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}

	*result = s
}

func loadEnvUint(key string, result *uint) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}

	n, err := strconv.Atoi(s)

	if err != nil {
		return
	}

	*result = uint(n)
}

type config struct {
	Listener listenerConfig `yaml:"listener" json:"listener"`
	DBConfig pgConfig       `yaml:"db" json:"db"`
	RedisConfig redisConfig `yaml:"redis" json:"redis"`
	JWTConfig jwtConfig `yaml:"jwt" json:"jwt"`
}

func (c *config) LoadFromEnv() {
	c.Listener.loadFromEnv()
	c.DBConfig.loadFromEnv()
	c.RedisConfig.loadFromEnv()
	c.JWTConfig.loadFromEnv()
}

func DefaultConfig() config {
	return config{
		Listener: defaultListenerConfig(),
		DBConfig: defaultPgConfig(),
		RedisConfig: defaultRedisConfig(),
		JWTConfig: defaultJwtConfig(),
	}
}

func loadConfigFromReader(r io.Reader, c *config) error {
	return yaml.NewDecoder(r).Decode(c)
}

func LoadConfigFromFile(fn string, c *config) error {
	_, err := os.Stat(fn)

	if err != nil {
		return err
	}

	f, err := os.Open(fn)

	if err != nil {
		return err
	}

	defer f.Close()

	return loadConfigFromReader(f, c)
}
