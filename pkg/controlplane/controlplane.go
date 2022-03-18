package controlplane

import (
	"sync"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/sirupsen/logrus"
)

var logger = common.Logger().WithFields(logrus.Fields{
	"pkg": "controlplane",
})

const (
	STOP_DEADLINE = 10
)

type Reconciler interface {
	Reconcile(stopCh <-chan struct{})
	Bootstrap()
	WaitStop() <-chan struct{}
}

type ControlPlane struct {
	recls []Reconciler

	mutex  sync.Mutex
	stopCh *chan struct{}
}

func New() *ControlPlane {
	return &ControlPlane{
		mutex: sync.Mutex{},
	}
}

func (c *ControlPlane) Start() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// only start once
	if c.stopCh == nil {
		stop := make(chan struct{})
		c.stopCh = &stop

		logger.Info("Bootstrapping...")
		// Bootstrap phase
		for i := range c.recls {
			c.recls[i].Bootstrap()
		}

		logger.Info("Reconciling...")
		// Reconcile phase
		for i := range c.recls {
			go c.recls[i].Reconcile(*c.stopCh)
		}
	}
}

func (c *ControlPlane) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.stopCh != nil {
		close(*c.stopCh)
		c.stopCh = nil

		// wait for all controller to stop?
		deadline := time.NewTimer(STOP_DEADLINE * time.Second)
		for i := range c.recls {
			select {
			case <-deadline.C:
				logger.Error("Deadline exceeds, Failed to wait all Controller to stopped, force stopping.")
			case <-c.recls[i].WaitStop():
			}
		}
		deadline.Stop()

		logger.Info("Stopped")
	}
}

func (c *ControlPlane) RegisterRecl(recl Reconciler) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.recls = append(c.recls, recl)

	// @TODO call new Reconciler if control plane has already started
}
