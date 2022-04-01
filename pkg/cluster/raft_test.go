package cluster_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/cluster"
	"github.com/hashicorp/raft"
)

type applier struct {
	kv map[string]string
}

func (a *applier) Apply(cmd *cluster.Command, log *raft.Log) *cluster.ApplyResponse {
	resp := &cluster.ApplyResponse{}
	fmt.Println("applier -> cmd: ", cmd)
	switch cmd.Op {
	case cluster.PutOp:
		a.kv[string(cmd.K)] = string(cmd.V)
	case cluster.DelOp:
		delete(a.kv, string(cmd.K))
	}

	return resp

}

func defaultRaConfig(kv map[string]string) *cluster.RaftConfig {
	return &cluster.RaftConfig{
		Ap: &applier{
			kv: kv,
		},
		SID:     "test-cluster-1",
		Addr:    "barry-x550la:7000",
		BaseDir: "./test",
	}
}

func bootRaft(kv map[string]string) (*cluster.RaftC, error) {
	cfg := defaultRaConfig(kv)

	if err := os.MkdirAll(cfg.BaseDir, 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	r, err := cluster.NewRaft(cfg)
	if err != nil {
		return r, err
	}

	future := r.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{{
			ID:      raft.ServerID(cfg.SID),
			Address: raft.ServerAddress(cfg.Addr),
		}},
	})

	if future.Error() != nil {
		return r, future.Error()
	}

	// wait for leader
	for {
		timeout := time.NewTimer(time.Second * 5)
		select {
		case isLead := <-r.LeaderCh():
			if isLead {
				timeout.Stop()
				return r, nil
			}
		case <-timeout.C:
			return r, fmt.Errorf("wait for leader timeout")

		}

	}
}

func destroyRaft(r *cluster.RaftC) {
	if r != nil {
		r.Shutdown().Error()
		os.RemoveAll(r.Cfg.BaseDir)
	}
}

func TestSingleRaft(t *testing.T) {
	kv := make(map[string]string)
	r, err := bootRaft(kv)
	defer destroyRaft(r)
	if err != nil {
		t.Error(err)
		return
	}

	var c cluster.Command
	c.Op = cluster.PutOp
	c.K = []byte("foo")
	c.V = []byte("123")

	b, err := json.Marshal(c)
	if err != nil {
		t.Error(err)
		return
	}

	future := r.Apply(b, 3*time.Second)
	if future.Error() != nil {
		t.Error(future.Error())
		return
	}

	_, ok := future.Response().(*cluster.ApplyResponse)
	if !ok {
		t.Errorf("invalid apply response")
		return
	}

	time.Sleep(5 * time.Second)

	fmt.Println(kv[string(c.K)])
	if kv[string(c.K)] == string(c.V) {
		return
	}

	t.Errorf("apply does not work properly")
}
