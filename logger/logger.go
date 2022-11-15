package logger

import (
	"time"

	adapter "github.com/axiomhq/axiom-go/adapters/zap"
	"github.com/axiomhq/axiom-go/axiom"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kamaesoft.visualstudio.com/kocto/_git/kocto/config"
)

func InitLogger(cfg config.Config) (*zap.SugaredLogger, error) {
	if cfg.Env == "development" {
		logger, err := devLogger(cfg.Log)
		if err != nil {
			return nil, err
		}
		return logger.Sugar(), nil
	}

	logger, err := prodLogger(cfg.Log)
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func devLogger(cfg config.LogConfig) (*zap.Logger, error) {
	logConfig := zap.NewDevelopmentConfig()
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return logConfig.Build()
}

func prodLogger(cfg config.LogConfig) (*zap.Logger, error) {
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

	return logger, nil
}

func syncer(logger *zap.Logger) {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		logger.Sync()
	}
}
