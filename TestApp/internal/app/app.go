package app

import (
	"TestApp/internal/config"
	"TestApp/internal/user"
	"TestApp/internal/user/db"
	"TestApp/pkg/client/postgresql"
	"TestApp/pkg/logging"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type App struct {
	cfg *config.Config

	router     *httprouter.Router
	httpServer *http.Server
	//grpcServer *grpc.Server

	userService user.Service
}

func NewApp(ctx context.Context, cfg *config.Config) (App, error) {
	logger := logging.GetLogger()

	pgClient, err := postgresql.NewClient(ctx, 3, cfg.Storage)
	if err != nil {
		return App{}, err
	}

	userStorage := db.NewRepository(pgClient, logger)
	userService := user.NewService(userStorage, logger)

	router := httprouter.New()
	userHandler := user.NewHandler(userService, logger)
	userHandler.Register(router)

	return App{
		cfg:         cfg,
		router:      router,
		userService: *userService,
		httpServer: &http.Server{
			Handler: router,
		},
	}, nil
}

func (a *App) Run() {
	logger := logging.GetLogger()
	logger.Info("starting server")

	var listener net.Listener
	var listenErr error

	if a.cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(os.Args[0])
		if err != nil {
			logger.Fatal(err)
		}

		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path: %s", socketPath)

		logger.Info("listening unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
	} else {
		logger.Info("Listening tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:           a.router,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	log.Fatal(server.Serve(listener))
}
