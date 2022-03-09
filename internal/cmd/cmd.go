package cmd

import (
	"path/filepath"

	"github.com/barrydevp/transcoorditor/pkg/api"
	"github.com/barrydevp/transcoorditor/pkg/api/controller"
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/service"
	"github.com/barrydevp/transcoorditor/pkg/store"
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
	s, err := store.NewMongoDBStore()
	if err != nil {
		panic(err)
	}

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
