package logic

import (
	"context"
	"sync"

	"github.com/Bazhenator/generator/configs"
	pb "github.com/Bazhenator/buffer/pkg/api/grpc"
	"github.com/Bazhenator/generator/pkg/connections/buffer"
	"github.com/Bazhenator/tools/src/logger"
	"golang.org/x/exp/rand"
)

type Service struct {
	c *configs.Config
	l *logger.Logger

	buffer buffer.Connection
}

func NewService(c *configs.Config, l *logger.Logger, con buffer.Connection) *Service {
	return &Service{
		c: c,
		l: l,

		buffer: con,
	}
}

func (s *Service) GenerateRequests(ctx context.Context, amount uint64) error {
	var wg sync.WaitGroup

	for i := uint64(0); i < amount; i++ {
		wg.Add(1)

		go func(id uint64) {
			defer wg.Done()
			
			req, err := generateSingleRequest(ctx, id, s.buffer.Client)
			if err != nil {
				s.l.Warn("buffer declined request", logger.NewField("req", req))
				s.l.Warn("failed to append request to buffer", logger.NewErrorField(err))
			}

			s.l.Debug("request appended successfully")
		}(i + 1)
	}

	wg.Wait()
	return nil
}

func generateSingleRequest(ctx context.Context, id uint64, client pb.BufferServiceClient) (*pb.Request, error) {
	req := &pb.Request{
		Id:           id,
		ClientId:     uint64(rand.Intn(1000)),
		Priority:     uint32(rand.Intn(10) + 1),
		CleaningType: uint32(rand.Intn(3) + 1),
	}

	_, err := client.AppendRequest(ctx, &pb.AppendRequestIn{Req: req})
	return req, err
}
