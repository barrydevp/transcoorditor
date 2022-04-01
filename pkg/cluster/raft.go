package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"go.etcd.io/bbolt"
)

const (
	raftAddr         = "127.0.0.1:7000"
	tcpTimeout       = 10 * time.Second
	tcpMaxPool       = 3
	snapshotRetain   = 10
	raftBaseDir      = ""
	raftLogFile      = "raft.db"
	raftLogCacheSize = 512
)

var (
	ErrNoApplier = errors.New("no applier")
)

type RaftConfig struct {
	SID      string
	Ap       Applier
	Addr     string
	DBFile   string
	BaseDir  string
	LogLevel string
}

func (rcfg *RaftConfig) raftBoltOptions() raftboltdb.Options {
	bopts := &bbolt.Options{}
	opts := raftboltdb.Options{
		BoltOptions: bopts,
		Path:        raftLogFile,
	}

	logFile := raftLogFile
	if rcfg.DBFile != "" {
		logFile = rcfg.DBFile
	}
	opts.Path = filepath.Join(rcfg.dir(), logFile)

	return opts
}

func (rcfg *RaftConfig) raftConfig() *raft.Config {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(rcfg.SID)
	config.LogLevel = rcfg.LogLevel
	// config.TrailingLogs = 5
	// config.SnapshotInterval = 20 * time.Second
	// config.SnapshotThreshold = 3

	return config
}

func (rcfg *RaftConfig) fsm() *fsm {
	return &fsm{
		ap: rcfg.Ap,
	}
}

func (rcfg *RaftConfig) raftAddr() string {
	if rcfg.Addr == "" {
		return raftAddr
	}

	return rcfg.Addr
}

func (rcfg *RaftConfig) dir() string {
	if rcfg.BaseDir == "" {
		return raftBaseDir
	}
	return rcfg.BaseDir
}

type RaftC struct {
	*raft.Raft
	Cfg *RaftConfig

	transport *raft.NetworkTransport
	snapshots *raft.FileSnapshotStore
	stable    *raftboltdb.BoltStore
	logs      *raft.LogCache
}

func NewRaft(cfg *RaftConfig) (*RaftC, error) {
	r := &RaftC{
		Cfg: cfg,
	}

	raftAddr := cfg.raftAddr()
	addr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return nil, fmt.Errorf("resolve tcp addr err: %w", err)
	}

	if r.transport, err = raft.NewTCPTransport(raftAddr, addr, tcpMaxPool, tcpTimeout, logger.Writer()); err != nil {
		return nil, fmt.Errorf("init transport err: %w", err)
	}

	// Create the snapshot store. This allows the Raft to truncate the log.
	if r.snapshots, err = raft.NewFileSnapshotStore(cfg.dir(), snapshotRetain, logger.Writer()); err != nil {
		return nil, fmt.Errorf("init snapshots err: %w", err)
	}

	if r.stable, err = raftboltdb.New(cfg.raftBoltOptions()); err != nil {
		return nil, fmt.Errorf("init raft boltdb err: %w", err)
	}

	if r.logs, err = raft.NewLogCache(raftLogCacheSize, r.stable); err != nil {
		return nil, fmt.Errorf("init log cache err: %v", err)
	}

	ra, err := raft.NewRaft(cfg.raftConfig(), cfg.fsm(), r.logs, r.stable, r.snapshots, r.transport)
	if err != nil {
		return nil, fmt.Errorf("init raft err: %w", err)
	}

	r.Raft = ra

	return r, nil

}

type fsm struct {
	ap Applier
}

func (f *fsm) Apply(log *raft.Log) interface{} {
	var cmd Command
	if err := json.Unmarshal(log.Data, &cmd); err != nil {
		panic(fmt.Errorf("malformed OpLog: %w", err))
	}
	logger.Info(" + [fsm] Apply Log: ", cmd.Op)

	var response *ApplyResponse
	if f.ap == nil {
		response = &ApplyResponse{
			Err: ErrNoApplier,
			Res: nil,
		}
	} else {
		response = f.ap.Apply(&cmd, log)

	}

	return response
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	fmt.Println("[fsm] Snapshot")
	return nil, nil
	// return &fsmSnapshot{}, nil
}

func (f *fsm) Restore(snapshot io.ReadCloser) error {
	fmt.Println("[fsm] Restore")
	o := make(map[string]string)
	if err := json.NewDecoder(snapshot).Decode(&o); err != nil {
		return err
	}

	// f.kv = o
	return nil
}

type fsmSnapshot struct {
}
