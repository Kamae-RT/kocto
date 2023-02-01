package kocto

import (
	"github.com/caarlos0/env/v7"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

const DefaultPort = "4444"

type Config struct {
	Env    Environment `env:"ENV"`
	Port   string      `env:"PORT"`
	Log    LogConfig
	DB     DBConfig
	Rabbit RabbitConfig
}

type LogConfig struct {
	Name    string `env:"LOG_NAME"`
	Token   string `env:"AXIOM_TOKEN"`
	Org     string `env:"AXIOM_ORG_ID"`
	Dataset string `env:"AXIOM_DATASET"`
}

type DBConfig struct {
	URL  string `env:"DATABASE_URL"`
	Name string `env:"DATABASE_NAME"`
}

type RabbitConfig struct {
	URL string `env:"RABBIT_URL"`
}

// Loads base configuration from the environment
func LoadConfig() (Config, error) {
	var cfg Config
	err := LoadInConfig(&cfg)

	return cfg, err
}

// LoadInConfig parses a struct containing env vars
// Use this function if you have your own extended config
func LoadInConfig(cfg any) error {
	return env.Parse(cfg)
}
