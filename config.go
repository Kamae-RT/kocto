package kocto

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

const DefaultPort = "4444"

type Config struct {
	Env    Environment
	Port   string
	Log    LogConfig
	DB     DBConfig
	Rabbit RabbitConfig
}

type LogConfig struct {
	Name    string
	Token   string
	Org     string
	Dataset string
}

type DBConfig struct {
	URL  string
	Name string
}

type RabbitConfig struct {
	URL string
}
