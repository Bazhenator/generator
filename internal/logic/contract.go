package logic

import "context"

type GeneratorService interface {
	GenerateRequests(context.Context, uint64) error 
}