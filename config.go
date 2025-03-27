package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
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

type pgConfig struct {
	Host string `yaml:"host" json:"host"`
	Port uint   `yaml:"port" json:"port"`

	DBName  string `yaml:"db_name" json:"db_name"`
	SslMode string `yaml:"ssl_mode" json:"ssl_mode"`
}

func (p pgConfig) ConnStr() string {
	return fmt.Sprintf("host=%s port=%d database=%s sslmode=%s", p.Host, p.Port, p.DBName, p.SslMode)
}

func defaultPgConfig() pgConfig {
	return pgConfig{
		Host:    "localhost",
		Port:    5432,
		DBName:  "nexmedis",
		SslMode: "disable",
	}
}

func (p *pgConfig) loadFromEnv() {
	loadEnvStr("DB_HOST", &p.Host)
	loadEnvUint("DB_PORT", &p.Port)
	loadEnvStr("DB_NAME", &p.DBName)
	loadEnvStr("DB_SSL", &p.SslMode)

}

type hashParam struct {
    Memory      uint32 `yaml:"memory" json:"memory"`
    Iterations  uint32 `yaml:"iterations" json:"iterations"`
    Parallelism uint8  `yaml:"parallelism" json:"parallelism"`
    SaltLength  uint32 `yaml:"salt_length" json:"salt_length"`
    KeyLength   uint32 `yaml:"key_length" json:"key_length"`
}

func defaultHashParam() hashParam {
    return hashParam{
        Memory:      64 * 1024,
        Iterations:  3,
        Parallelism: 2,
        SaltLength:  16,
        KeyLength:   32,
    }
}

func (h *hashParam) loadFromEnv() {
    loadEnvUint32("HASH_MEMORY", &h.Memory)
    loadEnvUint32("HASH_ITERATIONS", &h.Iterations)
    loadEnvUint8("HASH_PARALLELISM", &h.Parallelism)
    loadEnvUint32("HASH_SALT_LENGTH", &h.SaltLength)
    loadEnvUint32("HASH_KEY_LENGTH", &h.KeyLength)
}

func loadEnvUint32(key string, dest *uint32) {
    if val := os.Getenv(key); val != "" {
        if parsed, err := strconv.ParseUint(val, 10, 32); err == nil {
            *dest = uint32(parsed)
        }
    }
}

func loadEnvUint8(key string, dest *uint8) {
    if val := os.Getenv(key); val != "" {
        if parsed, err := strconv.ParseUint(val, 10, 8); err == nil {
            *dest = uint8(parsed)
        }
    }
}


type listenerConfig struct {
	Host string `yaml:"host" json:"host"`
	Port uint   `yaml:"port" json:"port"`
}

func (l listenerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

func defaultListenerConfig() listenerConfig {
	return listenerConfig{
		Host: "127.0.0.1",
		Port: 8080,
	}
}

func (l *listenerConfig) loadFromEnv() {
	loadEnvStr("LISTENER_HOST", &l.Host)
	loadEnvUint("LISTENER_PORT", &l.Port)
}


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
		Port: 8080,
		Password: "",
	}
}

func (l *redisConfig) loadFromEnv() {
	loadEnvStr("REDIS_HOST", &l.Host)
	loadEnvUint("REDIS_PORT", &l.Port)
	loadEnvStr("REDIS_PASSWORD", &l.Password)
}

type config struct {
	Listener listenerConfig `yaml:"listener" json:"listener"`
	DBConfig pgConfig       `yaml:"db" json:"db"`
	RedisConfig redisConfig `yaml:"redis" json:"redis"`
	HashParam hashParam `yaml:"hash" json:"hash"`
}

func (c *config) loadFromEnv() {
	c.Listener.loadFromEnv()
	c.DBConfig.loadFromEnv()
	c.RedisConfig.loadFromEnv()
	c.HashParam.loadFromEnv()
}

func defaultConfig() config {
	return config{
		Listener: defaultListenerConfig(),
		DBConfig: defaultPgConfig(),
		RedisConfig: defaultRedisConfig(),
		HashParam: defaultHashParam(),
	}
}

func loadConfigFromReader(r io.Reader, c *config) error {
	return yaml.NewDecoder(r).Decode(c)
}

func loadConfigFromFile(fn string, c *config) error {
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

func loadConfig(fn string) config {
	cfg := defaultConfig()
	cfg.loadFromEnv()
	if (len(fn) > 0) {
		err := loadConfigFromFile(fn, &cfg)
		if err != nil {
			log.Warn().Str("file", fn).Err(err)
		}
	}
	log.Debug().Any("config", cfg).Msg("config loaded")
	return cfg
}
