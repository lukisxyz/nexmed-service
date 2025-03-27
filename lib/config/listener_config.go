package config

import "fmt"

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