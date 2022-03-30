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
}

func New() *Cluster {
	return &Cluster{}
}

func (c *Cluster) Run(applier Applier) (err error) {
	rcfg := &RaftConfig{
		Ap:     applier,
		Addr:   viper.GetString("NODE_ADDR"),
		SID:    viper.GetString("NODE_ID"),
		DBFile: viper.GetString("RAFT_DB"),
	}

	if c.Ra, err = NewRaft(rcfg); err != nil {
		return err
	}

	return nil
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

    logger.Info(servers)

	future := c.Ra.BootstrapCluster(raft.Configuration{
		Servers: servers,
	})

	if future.Error() != nil {
		return future.Error()
	}

	return nil
}

func (c *Cluster) Join(node *Node) error {
	if c.Ra == nil {
		return ErrClusterNotRunning
	}

	future := c.Ra.AddVoter(node.id(), node.addr(), 0, 5*time.Second)

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
