package service

import (
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/sirupsen/logrus"
)

type Service struct {
	s store.Interface
	l *logrus.Entry
}

func NewService(s store.Interface) *Service {
	return &Service{
		s: s,
		l: common.Logger().WithFields(logrus.Fields{
			"pkg": "service",
		}),
	}
}
