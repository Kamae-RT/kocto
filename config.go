package kocto

import (
	"flag"
	"os"
	"strconv"
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

	LoadConfigWithoutParse(&cfg)

	flag.Parse()

	return cfg
}

func LoadConfigWithoutParse(cfg *Config) {
	env := ""
	flag.StringVar(&env, "env", GetEnvVar("ENV", "development"), "Environment (development|production)")
	cfg.Env = Environment(env)

	flag.StringVar(&cfg.Port, "port", GetEnvVar("PORT", "8080"), "Port the http server is running on")
	flag.StringVar(&cfg.DB.URL, "db-url", GetEnvVar("DATABASE_URL", "mongodb://localhost:27017/"), "MongoDB connection url")
	flag.StringVar(&cfg.DB.Name, "db-name", GetEnvVar("DATABASE_NAME", ""), "MongoDB connection url")
	flag.StringVar(&cfg.Rabbit.URL, "rabbit-url", GetEnvVar("RABBIT_URL", "amqp://localhost"), "RabbitMQ cluster url")
	flag.StringVar(&cfg.Log.Name, "log-name", GetEnvVar("LOG_NAME", ""), "")
	flag.StringVar(&cfg.Log.Org, "axiom-org", GetEnvVar("AXIOM_ORG_ID", ""), "")
	flag.StringVar(&cfg.Log.Token, "axiom-token", GetEnvVar("AXIOM_TOKEN", ""), "")
	flag.StringVar(&cfg.Log.Dataset, "axiom-dataset", GetEnvVar("AXIOM_DATASET", ""), "")
}

func GetEnvVar(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func GetBoolEnvVar(key string, defaultValue bool) bool {
	val := GetEnvVar(key, strconv.FormatBool(defaultValue))

	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return b
}
