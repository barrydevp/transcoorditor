package cluster

import (
	"github.com/hashicorp/raft"
)

type Applier interface {
	GetLastIndex()
	GetLastTerm()
}

type fsm struct {
	ra *raft.Raft
}

type fsmSnapshot struct {
}
