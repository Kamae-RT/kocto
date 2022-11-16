package kocto

import (
	"time"

	adapter "github.com/axiomhq/axiom-go/adapters/zap"
	"github.com/axiomhq/axiom-go/axiom"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a simple Wrapper for `*zap.SugaredLogger`
type Logger = *zap.SugaredLogger

// InitLogger setups the logger based on the Config environment
//
// in development:
//      - prints to console with colors
//      - default level is `debug`
// in production:
//      - prints to console in json format 
//      - batchs to Axiom.co
//      - default level is `info`
func InitLogger(cfg Config) (Logger, error) {
	if cfg.Env == "development" {
		return devLogger(cfg.Log)
	}

	return prodLogger(cfg.Log)
}

func devLogger(cfg LogConfig) (Logger, error) {
	logConfig := zap.NewDevelopmentConfig()
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

    logger, err := logConfig.Build()
    if err != nil {
        return nil, err
    }

    return logger.Sugar(), nil
}

func prodLogger(cfg LogConfig) (Logger, error) {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	consoleLogger, err := logConfig.Build()
	if err != nil {
		return nil, err
	}
	axiomCore, err := adapter.New(
		adapter.SetClientOptions(axiom.SetCloudConfig(
			cfg.Token,
			cfg.Org,
		)),
		adapter.SetDataset(cfg.Dataset),
		adapter.SetLevelEnabler(zapcore.InfoLevel),
	)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewTee(
		consoleLogger.Core(),
		axiomCore,
	)

	logger := zap.New(core).Named(cfg.Name)

	go syncer(logger)

    return logger.Sugar(), nil
}

func syncer(logger *zap.Logger) {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		logger.Sync()
	}
}
