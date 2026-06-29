package obs

import (
	"fmt"

	"github.com/razedwell/traxex/shared/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LoggerJsonFormat = "json"
)

func InitLogger(cfg config.LoggingConfig) (*zap.Logger, error) {
	zaplvl, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("Error parsing log level to zap: %w", err)
	}

	var zapcfg zap.Config
	if cfg.Format == LoggerJsonFormat {
		zapcfg = zap.NewProductionConfig()
	} else {
		zapcfg = zap.NewDevelopmentConfig()
	}
	zapcfg.Level = zap.NewAtomicLevelAt(zaplvl)

	logger, err := zapcfg.Build()
	if err != nil {
		return nil, fmt.Errorf("Error building zap logger from config: %w", err)
	}
	return logger, nil
}
