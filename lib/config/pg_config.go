package config

import "fmt"

type pgConfig struct {
	Host string `yaml:"host" json:"host"`
	Port uint   `yaml:"port" json:"port"`

	DBName  string `yaml:"db_name" json:"db_name"`
	SslMode string `yaml:"ssl_mode" json:"ssl_mode"`

	User string `yaml:"user" json:"user"`
	Secret string `yaml:"secret" json:"secret"`
}

func (p pgConfig) ConnStr() string {
	return fmt.Sprintf("host=%s port=%d database=%s sslmode=%s user=%s password=%s", p.Host, p.Port, p.DBName, p.SslMode, p.User, p.Secret)
}

func defaultPgConfig() pgConfig {
	return pgConfig{
		Host:    "localhost",
		Port:    5432,
		DBName:  "nexmedis",
		SslMode: "disable",
		User: "fahmi",
		Secret: "wap12345",
	}
}

func (p *pgConfig) loadFromEnv() {
	loadEnvStr("DB_HOST", &p.Host)
	loadEnvUint("DB_PORT", &p.Port)
	loadEnvStr("DB_NAME", &p.DBName)
	loadEnvStr("DB_SSL", &p.SslMode)
	loadEnvStr("DB_USER", &p.User)
	loadEnvStr("DB_SECRET", &p.Secret)

}