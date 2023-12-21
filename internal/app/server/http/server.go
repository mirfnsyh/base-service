package http

import (
	"context"
	"fmt"
	"github.com/mirfnsyh/base-service/internal/app/commons"
	"github.com/mirfnsyh/base-service/internal/app/controllers"
	services "github.com/mirfnsyh/base-service/internal/app/usecases"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
)

type IServer interface {
	StartApp()
}

type server struct {
	opt      commons.Options
	services *services.Services
}

func NewServer(opt commons.Options, services *services.Services) IServer {
	return &server{
		opt:      opt,
		services: services,
	}
}

func (s *server) StartApp() {
	var srv http.Server
	idleConnectionClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		logrus.Infoln("[API] Server is shutting down")

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logrus.Infof("[API] Fail to shutting down: %v", err)
		}
		close(idleConnectionClosed)
	}()

	srv.Addr = fmt.Sprintf("%s:%d", s.opt.AppCtx.GetAppOption().Host, s.opt.AppCtx.GetAppOption().Port)
	hOpt := controllers.ControllerOption{
		Options:  s.opt,
		Services: s.services,
	}
	srv.Handler = Router(hOpt)

	logrus.Infof("[API] HTTP serve at %s\n", srv.Addr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logrus.Infof("[API] Fail to start listen and server: %v", err)
	}

	<-idleConnectionClosed
	logrus.Infoln("[API] Bye")
}
