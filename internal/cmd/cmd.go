package cmd

import (
	"path/filepath"

	"github.com/barrydevp/transcoorditor/pkg/api"
	"github.com/barrydevp/transcoorditor/pkg/api/controller"
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/service"
	"github.com/barrydevp/transcoorditor/pkg/store/exclusive"
	"github.com/barrydevp/transcoorditor/pkg/store/mongodb"
)

var envFile = ".env"

func ApiServer() {
	// Loading env into viper config
	common.InitEnv(filepath.Join("./", envFile))
	common.InitLogger()

	// init api server
	apiSrv := api.NewServer()

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

	// add controler
	ctrl := controller.NewController(ac)
	// register routes
	ctrl.PublicRoutes(apiSrv.Srv)

	// Run server
	apiSrv.Run()

	// cleanup
	s.Close()

}
