package logic

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	pb "github.com/Bazhenator/buffer/pkg/api/grpc"
	"github.com/Bazhenator/generator/configs"
	"github.com/Bazhenator/generator/internal/logic/dto"
	"github.com/Bazhenator/generator/pkg/connections/buffer"
	"github.com/Bazhenator/tools/src/logger"
	"golang.org/x/exp/rand"
)

var (
	statusDeclined = uint32(0)
	statusAccepted = uint32(1)
)

type Service struct {
	c *configs.Config
	l *logger.Logger

	requests []*dto.Request
	buffer   buffer.Connection
}

func NewService(c *configs.Config, l *logger.Logger, con buffer.Connection) *Service {
	r := make([]*dto.Request, 0, c.RequestsAmount)

	return &Service{
		c: c,
		l: l,

		requests: r,
		buffer:   con,
	}
}

// GenerateRequests is a main logic func for generator service. It generates requests for cleaning service
// and tries to put it into buffer.
func (s *Service) GenerateRequests(ctx context.Context, amount uint64) error {
	var wg sync.WaitGroup

	requestsPerGenerator := amount / s.c.GensAmount
	extraRequests := uint64(0)

	for genID := uint64(0); genID < s.c.GensAmount; genID++ {
		wg.Add(1)

		if genID == s.c.GensAmount-1 {
			extraRequests = amount % s.c.GensAmount
		}

		go func(generatorID uint64, requests uint64) {
			defer wg.Done()
			s.l.Debug("Generator started", logger.NewField("generator_id", generatorID))

			for i := uint64(0); i < requests; i++ {
				req, err := generateSingleRequest(ctx, uint64(generatorID)*requests+i+1, generatorID, s.buffer.Client)
				if err == nil {
					req.Status = &statusAccepted
				}
				time.Sleep(time.Second * 2)

				collect := func() {
					s.requests = append(s.requests, &dto.Request{
						Id:           req.Id,
						ClientId:     req.ClientId,
						Priority:     req.Priority,
						CleaningType: req.CleaningType,
						GeneratorId:  *req.GeneratorId,
						Status:       *req.Status,
					})
				}

				if err != nil {
					s.l.Warn("buffer declined request", logger.NewField("req", req))
					collect()
					continue
				}

				collect()
				s.l.Debug("request appended successfully", logger.NewField("req", req))
			}

			s.l.Debug("Generator finished", logger.NewField("generator_id", generatorID))
		}(genID, requestsPerGenerator+extraRequests)
	}

	wg.Wait()

	return nil
}

// CreateReport creates report about generators' work statistics.
func (s *Service) CreateReport() error {
	reportFile, err := os.Create("report.txt")
	if err != nil {
		s.l.Error("failed to create report file", logger.NewErrorField(err))
		return err
	}
	defer reportFile.Close()

	var accepted, declined uint32

	for _, req := range s.requests {
		reportFile.WriteString("{ req:\n ")
		reportFile.WriteString("     { id:             " + fmt.Sprintf("%d", req.Id) + "\n")
		reportFile.WriteString("       client_id:      " + fmt.Sprintf("%d", req.ClientId) + "\n")
		reportFile.WriteString("       priority:       " + fmt.Sprintf("%d", req.Priority) + "\n")
		reportFile.WriteString("       cleaning_type:  " + fmt.Sprintf("%d", req.CleaningType) + "\n")
		reportFile.WriteString("       generator_id:   " + fmt.Sprintf("%d", req.GeneratorId)+ "\n")
		reportFile.WriteString("       status:         " + fmt.Sprintf("%d", req.Status)+ "\n")
		reportFile.WriteString("     }"+ "\n")
		reportFile.WriteString("}"+ "\n")

		if req.Status == statusDeclined {
			declined += 1
			continue
		}

		accepted += 1
	}

	reportFile.WriteString("Summary:\n")
	reportFile.WriteString("Accepted:" + fmt.Sprintf("%f", float64(accepted) / float64(len(s.requests))  * 100) + "%\n")
	reportFile.WriteString("Declined:" + fmt.Sprintf("%f", float64(declined) / float64(len(s.requests))  * 100) + "%\n")

	return nil
}

// generateSingleRequest is a supplimentary func for putting request into buffer of buffer service.
func generateSingleRequest(ctx context.Context, id uint64, generatorId uint64, client pb.BufferServiceClient) (*pb.Request, error) {
	req := &pb.Request{
		Id:           id,
		ClientId:     uint64(rand.Intn(1000)),
		Priority:     uint32(rand.Intn(5)),
		CleaningType: uint32(rand.Intn(10) + 1),
		GeneratorId:  &generatorId,
		Status:       &statusDeclined,
	}

	_, err := client.AppendRequest(ctx, &pb.AppendRequestIn{Req: req})
	return req, err
}
