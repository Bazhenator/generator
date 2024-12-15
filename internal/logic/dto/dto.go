package dto

type Request struct {
	Id           uint64
	ClientId     uint64
	CleaningType uint32
	Priority     uint32
	GeneratorId  uint64
	Status       uint32
}