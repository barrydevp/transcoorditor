package cluster

import (
	"encoding/json"

	"github.com/hashicorp/raft"
)

type Applier interface {
	// GetLastIndex()
	// GetLastTerm()
	Apply(c *Command, log *raft.Log) *ApplyResponse
}

type ApplyResponse struct {
	Err error
	Res interface{}
}

type CommandOp string

const (
	PutOp CommandOp = "put"
	DelOp CommandOp = "del"
	RpcOp CommandOp = "rpc"
)

type Command struct {
	Ns string
	Op CommandOp
	K  []byte
	V  []byte
}

func (c *Command) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Command) Decode(b []byte) error {
	return json.Unmarshal(b, c)
}

func NewRpcCmd(namespace string, method string, args ...interface{}) (*Command, error) {
	argByteArr := make([][]byte, len(args))
	for i, arg := range args {
		b, err := json.Marshal(arg)
		if err != nil {
			return nil, err
		}
		argByteArr[i] = b
	}

	argsByte, err := json.Marshal(argByteArr)
	if err != nil {
		return nil, err
	}

	return &Command{
		Ns: namespace,
		Op: RpcOp,
		K:  []byte(method),
		V:  argsByte,
	}, nil
}

func ParseRpcCmd(c *Command, args ...interface{}) error {
	argsByte := c.V
	argByteArr := make([][]byte, len(args))
	if err := json.Unmarshal(argsByte, &argByteArr); err != nil {
		return err
	}

	for i, arg := range args {
		if err := json.Unmarshal(argByteArr[i], arg); err != nil {
			return err
		}
	}

	return nil
}
