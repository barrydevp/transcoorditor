package cmd

import (
	"path/filepath"

	"github.com/barrydevp/transcoorditor/pkg/app"
	"github.com/barrydevp/transcoorditor/pkg/app/controller"
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/controlplane"
	"github.com/barrydevp/transcoorditor/pkg/service"
	"github.com/barrydevp/transcoorditor/pkg/store/exclusive"
	"github.com/barrydevp/transcoorditor/pkg/store/mongodb"
)

var envFile = ".env"

func RunApp() {
	// Loading env into viper config
	common.InitEnv(filepath.Join("./", envFile))
	common.InitLogger()

	// init api server
	apiSrv := app.NewServer()

	// init storage
	// s, err := store.NewMemoryStore()
	s, err := mongodb.NewStore()
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
