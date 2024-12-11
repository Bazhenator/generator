package configs

import (
	"errors"
	"os"

	"go.uber.org/multierr"

	"github.com/Bazhenator/tools/src/logger"
	grpcListener "github.com/Bazhenator/tools/src/server/grpc/listener"
)

const (
	envBufferService = "BUFFER_SERVICE"
)

// Config is a main configuration struct for application
type Config struct {
	Environment  string
	BufferHost   string
	Grpc         *grpcListener.GrpcConfig
	LoggerConfig *logger.LoggerConfig
}

// NewConfig returns application config instance
func NewConfig() (*Config, error) {
	var errorBuilder error

	grpcConfig, err := grpcListener.NewStandardGrpcConfig()
	multierr.AppendInto(&errorBuilder, err)

	loggerConfig, err := logger.NewLoggerConfig()
	multierr.AppendInto(&errorBuilder, err)

	bufferStr, ok := os.LookupEnv(envBufferService)
	if !ok {
		multierr.AppendInto(&errorBuilder, errors.New("BUFFER_SERVICE is not defined"))
	}

	if errorBuilder != nil {
		return nil, errorBuilder
	}

	glCfg := &Config{
		Grpc:         grpcConfig,
		LoggerConfig: loggerConfig,
		BufferHost:   bufferStr,
	}

	return glCfg, nil
}