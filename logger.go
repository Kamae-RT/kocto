package kocto

import (
	"time"

	adapter "github.com/axiomhq/axiom-go/adapters/zap"
	"github.com/axiomhq/axiom-go/axiom"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func InitLogger(cfg Config) (*Logger, error) {
	if cfg.Env == "development" {
		return devLogger(cfg.Log)
	}

	return prodLogger(cfg.Log)
}

func devLogger(cfg LogConfig) (*Logger, error) {
	logConfig := zap.NewDevelopmentConfig()
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

    logger, err := logConfig.Build()
    if err != nil {
        return nil, err
    }

    return &Logger{logger.Sugar()}, nil
}

func prodLogger(cfg LogConfig) (*Logger, error) {
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

    return &Logger{logger.Sugar()}, nil
}

func syncer(logger *zap.Logger) {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		logger.Sync()
	}
}
