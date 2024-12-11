package delivery

import (
	"context"

	"github.com/Bazhenator/generator/configs"
	"github.com/Bazhenator/generator/internal/logic"
	"github.com/Bazhenator/generator/pkg/api/grpc"
	"github.com/Bazhenator/tools/src/logger"
)

type GeneratorServer struct {
	generator.UnimplementedGeneratorServiceServer

	c *configs.Config
	l *logger.Logger

	logic logic.GeneratorService
}

func NewGeneratorServer(c *configs.Config, l *logger.Logger, logic logic.GeneratorService) *GeneratorServer {
	return &GeneratorServer{
		c: c,
		l: l,

		logic: logic,
	}
}

// StartGenerator starts generation of requests for cleaning service
func (s *GeneratorServer) StartGenerator(ctx context.Context, in *generator.StartGeneratorIn) (*generator.StartGeneratorOut, error) {
	s.l.InfoCtx(ctx, "StartGenerator called", logger.NewField("request_amount", in.RequestsAmount))

	err := s.logic.GenerateRequests(ctx, in.RequestsAmount)
	if err != nil {
		s.l.Error("Failed to generate requests", logger.NewErrorField(err))
		return &generator.StartGeneratorOut{Status: false}, err
	}

	return &generator.StartGeneratorOut{Status: true}, nil
}