package controlplane

import (
	"sync"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/cluster"
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
	c     *cluster.Cluster
	recls []Reconciler

	mutex      sync.Mutex
	stopCh     *chan struct{}
	reclStopCh *chan struct{}
}

func New(c *cluster.Cluster) *ControlPlane {
	return &ControlPlane{
		c:     c,
		mutex: sync.Mutex{},
	}
}

func (ctrl *ControlPlane) watchClusterLeader(leaderCh <-chan bool) {
	// do we need to check leadership before watch on leaderCh for manually start reconciler?
	logger.Info("Watch leadership...")

	for {
		select {
		case isLeader := <-leaderCh:
			ctrl.mutex.Lock()
			if isLeader {
				logger.Info("[+] on Leader")
				ctrl.unsafeStartReconciler()
			} else {
				logger.Info("[-] on Follower")
				ctrl.unsafeStopReconciler()
			}
			ctrl.mutex.Unlock()
		case <-*ctrl.stopCh:
			// the stopCh's sender must manully stop reconciling
			return
		}

	}
}

func (ctrl *ControlPlane) Run() error {
	ctrl.mutex.Lock()
	defer ctrl.mutex.Unlock()

	if ctrl.stopCh != nil {
		return nil
	}

	// watch leader transfer on cluster leader channel
	var leaderCh <-chan bool
	if ctrl.c != nil {
		ch, err := ctrl.c.LeaderCh()
		if err != nil {
			return err
		}
		leaderCh = ch
	} else {
		// for non-cluster mode
		// @FIXME: this channel must be close some where when stop
		ch := make(chan bool, 1)
		// simulate leadership
		ch <- true

		leaderCh = ch
	}

	stop := make(chan struct{})
	ctrl.stopCh = &stop
	logger.Info("Running...")

	go ctrl.watchClusterLeader(leaderCh)

	return nil
}

func (c *ControlPlane) unsafeStartReconciler() {
	// only start once
	if c.reclStopCh != nil {
		return
	}

	stop := make(chan struct{})
	c.reclStopCh = &stop

	logger.Info("+ Boot recl")
	// Bootstrap phase
	for i := range c.recls {
		c.recls[i].Bootstrap()
	}

	logger.Info("+ Reconcile recl")
	// Reconcile phase
	for i := range c.recls {
		go c.recls[i].Reconcile(*c.reclStopCh)
	}
}

func (c *ControlPlane) unsafeStopReconciler() {
	if c.reclStopCh != nil {
		close(*c.reclStopCh)
		c.reclStopCh = nil

		// wait for all controller to stop? => turn deadline timeout as an option
		deadline := time.NewTimer(STOP_DEADLINE * time.Second)
		for i := range c.recls {
			select {
			case <-deadline.C:
				logger.Error("Deadline exceeds, Failed to wait all Controller to stopped, force stopping.")
			case <-c.recls[i].WaitStop():
			}
		}
		deadline.Stop()

		logger.Info("- recl Stopped")
	}
}

func (c *ControlPlane) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.stopCh != nil {
		close(*c.stopCh)
		c.stopCh = nil
	}

	c.unsafeStopReconciler()

	logger.Info("Stopped!")
}

func (c *ControlPlane) RegisterRecl(recl Reconciler) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.recls = append(c.recls, recl)

	// @TODO call new Reconciler if control plane has already started
}
