package action

import (
	"github.com/barrydevp/transcoorditor/pkg/store"
)

type Action struct {
	s store.Interface
}

func NewAction(s store.Interface) *Action {
	return &Action{
		s: s,
	}
}
