package controlplane_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/controlplane"
)

type reconciler struct {
	t testing.B
}

func (r *reconciler) Bootstrap() {
	fmt.Println("bootstrap has been called")
}

func (r *reconciler) Reconcile(stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			fmt.Println("stop received.")
			break
		case <-time.After(1 * time.Second):
			fmt.Println("tick.")
		}
	}
}

func (r *reconciler) WaitStop() <-chan struct{} {
	c := make(chan struct{})
	defer close(c)

	return c
}

func TestControlplane(t *testing.T) {
	recl := &reconciler{}
	ctrlplane := controlplane.New()

	ctrlplane.RegisterRecl(recl)

	ctrlplane.Start()

	time.Sleep(5 * time.Second)

	ctrlplane.Stop()
}
