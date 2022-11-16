package kocto

import (
	"flag"
	"os"
)

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

func LoadConfig() Config {
	var cfg Config

	env := ""
	flag.StringVar(&env, "env", getEnvVar("ENV", "development"), "Environment (development|production)")
	cfg.Env = Environment(env)

	flag.StringVar(&cfg.Port, "port", getEnvVar("PORT", "8080"), "Port the http server is running on")
	flag.StringVar(&cfg.DB.URL, "db-url", getEnvVar("DATABASE_URL", "mongodb://localhost:27017/"), "MongoDB connection url")
	flag.StringVar(&cfg.DB.Name, "db-name", getEnvVar("DATABASE_NAME", ""), "MongoDB connection url")
	flag.StringVar(&cfg.Rabbit.URL, "rabbit-url", getEnvVar("RABBIT_URL", "amqp://localhost"), "RabbitMQ cluster url")
	flag.StringVar(&cfg.Log.Name, "log-name", getEnvVar("LOG_NAME", ""), "")
	flag.StringVar(&cfg.Log.Org, "axiom-org", getEnvVar("AXIOM_ORG_ID", ""), "")
	flag.StringVar(&cfg.Log.Token, "axiom-token", getEnvVar("AXIOM_TOKEN", ""), "")

	return cfg
}

func getEnvVar(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}
