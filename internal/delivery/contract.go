package delivery

import "context"

type GeneratorService interface {
	GenerateRequests(context.Context, uint64) error 
	CreateReport() error
}