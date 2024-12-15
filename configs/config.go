package configs

import (
	"errors"
	"os"
	"strconv"

	"go.uber.org/multierr"

	"github.com/Bazhenator/tools/src/logger"
	grpcListener "github.com/Bazhenator/tools/src/server/grpc/listener"
)

const (
	EnvBufferService  = "BUFFER_SERVICE"
	EnvRequestsAmount = "REQUESTS_AMOUNT"
	EnvGensAmount     = "GENERATORS_AMOUNT"
)

// Config is a main configuration struct for application
type Config struct {
	Environment    string
	BufferHost     string
	GensAmount     uint64
	RequestsAmount uint64
	Grpc           *grpcListener.GrpcConfig
	LoggerConfig   *logger.LoggerConfig
}

// NewConfig returns application config instance
func NewConfig() (*Config, error) {
	var errorBuilder error

	grpcConfig, err := grpcListener.NewStandardGrpcConfig()
	multierr.AppendInto(&errorBuilder, err)

	loggerConfig, err := logger.NewLoggerConfig()
	multierr.AppendInto(&errorBuilder, err)

	bufferStr, ok := os.LookupEnv(EnvBufferService)
	if !ok {
		multierr.AppendInto(&errorBuilder, errors.New("BUFFER_SERVICE is not defined"))
	}

	EnvGensAmountStr, ok := os.LookupEnv(EnvGensAmount)
	if !ok {
		multierr.AppendInto(&errorBuilder, errors.New("GENERATORS_AMOUNT is not defined"))
	}

	gensAmount, err := strconv.Atoi(EnvGensAmountStr)
	multierr.AppendInto(&errorBuilder, err)

	EnvRequestsAmountStr, ok := os.LookupEnv(EnvRequestsAmount)
	if !ok {
		multierr.AppendInto(&errorBuilder, errors.New("REQUESTS_AMOUNT is not defined"))
	}

	reqsAmount, err := strconv.Atoi(EnvRequestsAmountStr)
	multierr.AppendInto(&errorBuilder, err)

	if errorBuilder != nil {
		return nil, errorBuilder
	}

	glCfg := &Config{
		Grpc:           grpcConfig,
		LoggerConfig:   loggerConfig,
		BufferHost:     bufferStr,
		GensAmount:     uint64(gensAmount),
		RequestsAmount: uint64(reqsAmount),
	}

	return glCfg, nil
}
