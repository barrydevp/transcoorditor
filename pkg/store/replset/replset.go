package replset

import (
	"errors"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/cluster"
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/sirupsen/logrus"
)

var (
	logger = common.Logger().WithFields(logrus.Fields{
		"pkg": "store/replset",
	})

	executeTimeout = 10 * time.Second

	ErrUnExpectedResponse   = errors.New("replset has return un-expected response")
	ErrRpcUnsupported       = errors.New("this rpc method does not supported")
	ErrNamespaceUnsupported = errors.New("this rpc namespace does not supported")
	ErrCmdUnsupported       = errors.New("this cmd does not supported")
)

type replsetBackend struct {
	*store.Backend
	s store.Interface
	c *cluster.Cluster

	internalSession     *sessionRepo
	internalParticipant *participantRepo
}

func NewReplStore(s store.Interface, c *cluster.Cluster) (*replsetBackend, error) {
	replset := &replsetBackend{
		s: s,
		c: c,
	}

	replset.internalSession = NewSession(replset)
	replset.internalParticipant = NewParticipant(replset)
	replset.Backend = &store.Backend{
		SessionImpl:     replset.internalSession,
		ParticipantImpl: replset.internalParticipant,
	}

	return replset, nil
}

// @overide
func (s *replsetBackend) Close() {
	s.s.Close()
}

func NewApplyErr(err error) *cluster.ApplyResponse {
	return &cluster.ApplyResponse{
		Err: err,
	}
}

func (s *replsetBackend) executeRPC(c *cluster.Command) *cluster.ApplyResponse {
	switch c.Ns {
	case "Session":
		return s.internalSession.executeRPC(c)
	case "Participant":
		return s.internalParticipant.executeRPC(c)
	}

	return NewApplyErr(ErrNamespaceUnsupported)
}

// @overide for Applier
func (s *replsetBackend) Apply(c *cluster.Command) *cluster.ApplyResponse {
	switch c.Op {
	case cluster.RpcOp:
		return s.executeRPC(c)
	}

	return NewApplyErr(ErrCmdUnsupported)
}
