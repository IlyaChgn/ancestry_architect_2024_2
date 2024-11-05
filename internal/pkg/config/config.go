package config

type RequestUUIDKey string
type LoggerKey string

const (
	OutputLogPath                        = "logs.json"
	ErrorOutputLogPath                   = "stderr err_logs.json"
	RequestUUIDContextKey RequestUUIDKey = "requestUUID"
	LoggerContextKey      LoggerKey      = "logger"
)

type AppServerConfig struct {
	Host    string   `yaml:"host"`
	Port    string   `yaml:"port"`
	Timeout int      `yaml:"timeout"`
	Origins []string `yaml:"origins"`
	Headers []string `yaml:"headers"`
	Methods []string `yaml:"methods"`
}

type PostgresConfig struct {
	Username string `env:"POSTGRES_USERNAME"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DBName   string `env:"DATABASE_NAME"`
}

type RedisConfig struct {
	Password string `env:"REDIS_PASSWORD"`
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
}

type AppRedisConfig struct {
	RedisConfig
	DB int `env:"APPLICATION_DB"`
}

type AppConfig struct {
	Server       AppServerConfig   `yaml:"server"`
	AdminService AdminServerConfig `yaml:"admin_service"`
	Postgres     PostgresConfig
	Redis        AppRedisConfig
}

type AdminServerConfig struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}

type AdminRedisConfig struct {
	RedisConfig
	DB int `env:"ADMIN_DB"`
}

type AdminConfig struct {
	Server   AdminServerConfig `yaml:"admin_service"`
	Postgres PostgresConfig
	Redis    AdminRedisConfig
}
