package cmd

import (
	"path/filepath"
	"strings"

	"github.com/barrydevp/transcoorditor/pkg/app"
	"github.com/barrydevp/transcoorditor/pkg/app/controller"
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/controlplane"
	"github.com/barrydevp/transcoorditor/pkg/service"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/barrydevp/transcoorditor/pkg/store/boltdb"
	"github.com/barrydevp/transcoorditor/pkg/store/exclusive"
	"github.com/barrydevp/transcoorditor/pkg/store/memory"
	"github.com/barrydevp/transcoorditor/pkg/store/mongodb"
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
		panic(err)
	}
	// add process synchronize exclusive locking when access data
	s, err = exclusive.NewStore(s)

	// init action
	ac := service.NewService(s)

	// init controlplane
	ctrlplane := controlplane.New()

	// add controler
	ctrl := controller.NewController(ac)
	// register routes
	ctrl.PublicRoutes(apiSrv.Srv)
	// register reconciler
	ctrl.RegisterReconciler(ctrlplane)

	// Run controlplane
	ctrlplane.Start()

	// Run server -> Start blocking from here
	apiSrv.Run()

	// cleanup
	ctrlplane.Stop()
	s.Close()

}
