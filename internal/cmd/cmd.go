package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/barrydevp/transcoorditor/pkg/app"
	"github.com/barrydevp/transcoorditor/pkg/app/controller"
	"github.com/barrydevp/transcoorditor/pkg/cluster"
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/controlplane"
	"github.com/barrydevp/transcoorditor/pkg/service"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/barrydevp/transcoorditor/pkg/store/boltdb"
	"github.com/barrydevp/transcoorditor/pkg/store/exclusive"
	"github.com/barrydevp/transcoorditor/pkg/store/memory"
	"github.com/barrydevp/transcoorditor/pkg/store/mongodb"
	"github.com/barrydevp/transcoorditor/pkg/store/replset"
	"github.com/spf13/viper"
)

var envFile = ".env"

func initStore() (store.Interface, error) {
	switch strings.ToLower(viper.GetString("BACKEND_STORE")) {
	case "mongodb":
		return mongodb.NewStore()
	case "memory":
		return memory.NewStore()
	default:
		return boltdb.NewStore()
	}
}

func RunApp() {
	// Loading env into viper config
	common.InitEnv(filepath.Join("./", envFile))
	common.InitLogger()

	// init api server
	apiSrv := app.NewServer()

	// init storage
	s, err := initStore()
	if err != nil {
		panic(fmt.Errorf("cannot init store: %w", err))
	}

	// add process synchronize exclusive locking when access data
	s, _ = exclusive.NewStore(s)

	// // init cluster
	clus := cluster.New()
	// replset store
	rsStore, _ := replset.NewReplStore(s, clus)

	if err := clus.Run(rsStore); err != nil {
		panic(fmt.Errorf("cannot run cluster: %w", err))
	}
	s = rsStore

	// init action
	ac := service.NewService(s)

	// init controlplane
	ctrlplane := controlplane.New(clus)

	// add controler
	ctrl := controller.NewController(clus, ac)
	// register routes
	ctrl.SystemRoutes(apiSrv.Srv)
	// register routes
	ctrl.PublicRoutes(apiSrv.Srv)
	// register reconciler
	ctrl.RegisterReconciler(ctrlplane)

	// Run controlplane
	if err = ctrlplane.Run(); err != nil {
		panic(fmt.Errorf("cannot run controlplane: %w", err))
	}

	// Run server -> Start blocking from here
	apiSrv.Run()

	// cleanup
	ctrlplane.Stop()
	clus.Stop()
	s.Close()

}
