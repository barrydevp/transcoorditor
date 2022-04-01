package exclusive

import (
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/hashicorp/raft"
)

type replsetRepo struct {
	*baseRepo
	s store.Replset
}

func NewReplset(s store.Replset) *replsetRepo {
	return &replsetRepo{
		baseRepo: newBaseRepo(),
		s:        s,
	}
}

func (s *replsetRepo) SaveLastLog(log *raft.Log) (err error) {
	s.withLock("replset_log", func() {
		err = s.s.SaveLastLog(log)
	})

	return
}

func (s *replsetRepo) GetLastLog() (log *raft.Log, err error) {
	s.withLock("replset_log", func() {
		log, err = s.s.GetLastLog()
	})

	return
}
