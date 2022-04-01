package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/store/mongodb"
	"github.com/hashicorp/raft"
)

func init() {
	common.InitEnv(".env")
}

func TestSaveLastLog(t *testing.T) {
	s, err := mongodb.NewStore()
	if err != nil {
		t.Error(err)
		return
	}

	err = s.Replset().SaveLastLog(&raft.Log{})
	if err != nil {
		t.Error(err)
	}
}

func TestGetLastLog(t *testing.T) {
	s, err := mongodb.NewStore()
	if err != nil {
		t.Error(err)
		return
	}

	log, err := s.Replset().GetLastLog()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(log)
}
