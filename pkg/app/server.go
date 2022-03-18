package app

import (
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ApiServer struct {
	Srv     *fiber.App
	CloseCh chan struct{}
	l       *logrus.Entry
}

func (s *ApiServer) WithGracefulShutdown() {
	s.CloseCh = make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := s.Srv.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			s.l.Error("Oops... Cannot shutdown Server! Reason: %v. Force shutdown!\n", err)
		}

		close(s.CloseCh)
	}()
}

func (s *ApiServer) BindCommonMiddleware() {
	// Logger middelware
	s.Srv.Use(logger.New(logger.Config{
		// Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
		// Output: os.Stdout,
	}))

}

// FiberConfig func for configuration Fiber app.
// See: https://docs.gofiber.io/api/fiber#config
func getFiberConfig() fiber.Config {
	// Define server settings.
	readTimeoutSecondsCount, _ := strconv.Atoi(viper.GetString("SERVER_READ_TIMEOUT"))

	// Return Fiber configuration.
	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}
}

func getServerUrl() string {
	portNumber := viper.GetString("PORT")

	if portNumber == "" {
		portNumber = "8000"
	}

	return ":" + portNumber
}

func NewServer() *ApiServer {
	// Make config
	config := getFiberConfig()

	// Define new Fiber app with config
	srv := fiber.New(config)

    server := &ApiServer{
		Srv:     srv,
		CloseCh: nil,
		l: common.Logger().WithFields(logrus.Fields{
			"pkg": "api/server",
		}),
	}

    // Common middlewares
    server.BindCommonMiddleware()

    return server
}

func (s *ApiServer) Run() {

	// graceful shutdown
	s.WithGracefulShutdown()

	// Run server.
	if err := s.Srv.Listen(getServerUrl()); err != nil {
		s.l.Error("Oops... Server is not running! Reason: %v", err)
	}

	if s.CloseCh != nil {
		<-s.CloseCh
	}

	s.l.Info("Server is shutdown.")
}
