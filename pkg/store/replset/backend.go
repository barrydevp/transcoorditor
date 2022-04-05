package replset

import (
	"errors"
	"fmt"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/cluster"
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/hashicorp/raft"
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
	ErrStaleCmd             = errors.New("stale command")
)

type replsetBackend struct {
	*store.Backend
	s store.Interface
	c *cluster.Cluster

	internalSession     *sessionRepo
	internalParticipant *participantRepo
	// indicate that the store is in replaying cmd state which is happend when starting replset server (early period after you run server in replset mode)
	replaying bool
	lastLog   *raft.Log
}

func NewReplStore(s store.Interface, c *cluster.Cluster) (*replsetBackend, error) {
	rs := &replsetBackend{
		s: s,
		c: c,
	}

	rs.internalSession = NewSession(rs)
	rs.internalParticipant = NewParticipant(rs)
	rs.Backend = &store.Backend{
		SessionImpl:     rs.internalSession,
		ParticipantImpl: rs.internalParticipant,
	}
	// retrive last log from backend store
	lastLog, err := s.Replset().GetLastLog()
	if err != nil {
		return rs, err
	}
	if lastLog != nil {
		logger.Info("Detect last log, change into replaying mode")
		rs.replaying = true
		rs.lastLog = lastLog
	}

	return rs, nil
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
func (rs *replsetBackend) Apply(c *cluster.Command, log *raft.Log) *cluster.ApplyResponse {
	// skip old command when in replaying mode
	if rs.replaying && rs.lastLog != nil {
		// @TODO: verify log.Term to determine this replset node is in in-consistent state
		if log.Index <= rs.lastLog.Index {
			logger.Debug("Skip old cmd, comming index: ", log.Index, " last persisted index: ", rs.lastLog.Index)
			return NewApplyErr(ErrStaleCmd)
		} else {
			// disable replaying mode
			rs.replaying = false
		}
	}

	var resp *cluster.ApplyResponse

	switch c.Op {
	case cluster.RpcOp:
		resp = rs.executeRPC(c)
	default:
		resp = NewApplyErr(ErrCmdUnsupported)
	}

	err := rs.s.Replset().SaveLastLog(log)
	if err != nil {
		return NewApplyErr(fmt.Errorf("cannot save last log: %w", err))
	}
	// rs.lastLog = log

	return resp
}
