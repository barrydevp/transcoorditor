package cluster

import (
	"errors"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	logger = common.Logger().WithFields(logrus.Fields{
		"pkg": "cluster",
	})

	ErrClusterNotRunning    = errors.New("cluster is not running")
	ErrUnknownApplyResponse = errors.New("unknown apply response")
	ErrNotLeader            = errors.New("not leader")
)

type Cluster struct {
	Ra      *RaftC
	running bool
}

type Node struct {
	ID   string
	Host string
}

func (n *Node) id() raft.ServerID {
	return raft.ServerID(n.ID)
}

func (n *Node) addr() raft.ServerAddress {
	return raft.ServerAddress(n.Host)
}

type ClusterRsConf struct {
	RsName string
	Nodes  []*Node
	Leader *Node
}

func New() *Cluster {
	return &Cluster{}
}

func (c *Cluster) Run(applier Applier) (err error) {
	rcfg := &RaftConfig{
		Ap:      applier,
		Addr:    viper.GetString("NODE_ADDR"),
		SID:     viper.GetString("NODE_ID"),
		DBFile:  viper.GetString("RAFT_DB"),
		BaseDir: viper.GetString("CLUSTER_BASE_DIR"),
	}

	if c.Ra, err = NewRaft(rcfg); err != nil {
		return err
	}

	return nil
}

func (c *Cluster) Stop() {
    // snapshot

}

func (c *Cluster) RsInitiate(rsconf *ClusterRsConf) error {
	if c.Ra == nil {
		return ErrClusterNotRunning
	}

	servers := make([]raft.Server, len(rsconf.Nodes))
	for i, node := range rsconf.Nodes {
		servers[i] = raft.Server{
			ID:      node.id(),
			Address: node.addr(),
		}
	}

	future := c.Ra.BootstrapCluster(raft.Configuration{
		Servers: servers,
	})

	if future.Error() != nil {
		return future.Error()
	}

	return nil
}

func (c *Cluster) Join(node *Node) error {
	if err := c.AssertRunning(); err != nil {
		return err
	}

	future := c.Ra.AddVoter(node.id(), node.addr(), 0, 5*time.Second)

	if future.Error() != nil {
		return future.Error()
	}

	return nil
}

func (c *Cluster) Left(node *Node) error {
	if err := c.AssertRunning(); err != nil {
		return err
	}

	future := c.Ra.RemoveServer(node.id(), 0, 5*time.Second)

	if future.Error() != nil {
		return future.Error()
	}

	return nil
}

func (c *Cluster) SID() string {
	if c.Ra != nil {
		return c.Ra.Cfg.SID
	}

	return ""
}

func (c *Cluster) Execute(cmd *Command, timeout time.Duration) (interface{}, error) {
	buf, err := cmd.Encode()
	if err != nil {
		return nil, err
	}

	future := c.Ra.Apply(buf, timeout)
	if future.Error() != nil {
		return nil, future.Error()
	}

	if resp, ok := future.Response().(*ApplyResponse); ok {
		if resp.Err != nil {
			return nil, resp.Err
		}
		return resp.Res, nil
	}

	return nil, ErrUnknownApplyResponse
}

func (c *Cluster) GetRaftConf() (*raft.Configuration, error) {
	if err := c.AssertRunning(); err != nil {
		return nil, err
	}

	confFuture := c.Ra.GetConfiguration()

	if confFuture.Error() != nil {
		return nil, confFuture.Error()
	}

	conf := confFuture.Configuration()

	return &conf, nil
}

func (c *Cluster) GetConf() (*ClusterRsConf, error) {
	raftConf, err := c.GetRaftConf()
	if err != nil {
		return nil, err
	}

	rsconf := &ClusterRsConf{}
	leaderAddr := c.Ra.Leader()
	for _, server := range raftConf.Servers {
		n := &Node{
			ID:   string(server.ID),
			Host: string(server.Address),
		}
		rsconf.Nodes = append(rsconf.Nodes, n)
		if server.Address == leaderAddr {
			rsconf.Leader = n
		}
	}

	return rsconf, nil
}

func (c *Cluster) AssertRunning() error {
	if c.Ra == nil {
		return ErrClusterNotRunning
	}
	return nil
}

func (c *Cluster) AssertLeader() error {
	// ensure cluster is running when call
	if c.Ra.State() != raft.Leader {
		return ErrNotLeader
	}
	return nil
}

func (c *Cluster) VerifyLeader() error {
	verifyFuture := c.Ra.VerifyLeader()
	if verifyFuture.Error() != nil {
		return verifyFuture.Error()
	}

	return nil
}

func (c *Cluster) Leader() (*Node, error) {
	conf, err := c.GetConf()
	if err != nil {
		return nil, err
	}

	return conf.Leader, nil
}

func (c *Cluster) LeaderHost() string {
	// ensure cluster is running when call
	return string(c.Ra.Leader())
}

func (c *Cluster) Stats() (map[string]string, error) {
	if err := c.AssertRunning(); err != nil {
		return nil, err
	}

	return c.Ra.Stats(), nil
}

func (c *Cluster) LeaderCh() (<-chan bool, error) {
	if err := c.AssertRunning(); err != nil {
		return nil, err
	}

	return c.Ra.LeaderCh(), nil
}
