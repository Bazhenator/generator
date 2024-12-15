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

var (
	statusDeclined = 0
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
	mu := sync.Mutex{}
  
	requestsPerGenerator := amount / s.c.GensAmount
	extraRequests := uint64(0)
  
	for genID := uint64(0); genID < s.c.GensAmount; genID++ {
	  wg.Add(1)

	  if genID == s.c.GensAmount - 1 {
		extraRequests = amount % s.c.GensAmount
	  }
  
	  go func(generatorID uint64, requests uint64) {
		defer wg.Done()
		s.l.Debug("Generator started", logger.NewField("generator_id", generatorID))
  
		for i := uint64(0); i < requests; i++ {
		  req, err := generateSingleRequest(ctx, uint64(generatorID)*requests+i+1, generatorID, s.buffer.Client)
		  if err != nil {
			s.l.Warn("buffer declined request", logger.NewField("req", req))
  
			continue
		  }
  
		  mu.Lock()
		  //req.Status = 1 // accepted
		  mu.Unlock()
  
		  s.l.Debug("request appended successfully", logger.NewField("req", req))
		}
  
		s.l.Debug("Generator finished", logger.NewField("generator_id", generatorID))
	  }(genID, requestsPerGenerator + extraRequests)
	}
  
	wg.Wait()
  
	return nil
  }

func generateSingleRequest(ctx context.Context, id uint64, generatorId uint64, client pb.BufferServiceClient) (*pb.Request, error) {
	req := &pb.Request{
	  Id:           id,
	  ClientId:     uint64(rand.Intn(1000)),
	  Priority:     uint32(generatorId),
	  CleaningType: uint32(rand.Intn(10) + 1),
	  //GeneratorId: generatorId,
	  //Status: statusDeclined,
	}

	_, err := client.AppendRequest(ctx, &pb.AppendRequestIn{Req: req})
	return req, err
}
