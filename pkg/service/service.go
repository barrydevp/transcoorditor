package service

import (
	"errors"

	"github.com/barrydevp/transcoorditor/pkg/store"
)

var (
	ErrNotFound        = errors.New("NOT_FOUND")
	ErrInvalidArgument = errors.New("INVALID_ARGUMENT")
	ErrPreconditionFailed    = errors.New("PRECONDITION_FAILED")
	ErrAborted         = errors.New("ABORTED")
)

type Service struct {
	s store.Interface
}

func NewService(s store.Interface) *Service {
	return &Service{
		s: s,
	}
}
