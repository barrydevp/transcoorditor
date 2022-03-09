package service

import (
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type Service struct {
	s store.Interface
}

func NewService(s store.Interface) *Service {
	return &Service{
		s: s,
	}
}
