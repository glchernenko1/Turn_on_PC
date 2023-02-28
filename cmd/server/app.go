package main

import (
	"Turn_on_PC/pkg/logging"
	"github.com/julienschmidt/httprouter"

	"Turn_on_PC/internal/server/DB/postgres"
	"Turn_on_PC/internal/server/config"
	"Turn_on_PC/internal/server/handlers"
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("Start Server")
	logger.Info("reading configuration")
	cfg := config.GetConfig()
	db := postgres.NewDB(cfg, logger)
	router := httprouter.New()
	handler := handlers.NewHandler(logger, db)
	handler.Register(router)
	handlerWS := handlers.NewHandlerWs(logger)
	handlerWS.Register(router)
	start(router, cfg)

}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Infoln("start application")

	var listenErr error
	var listener net.Listener

	logger.Infoln("listen tcp")
	listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Server.BindIP, cfg.Server.Port))
	logger.Infof("server is listening: http://%s:%s ", cfg.Server.BindIP, cfg.Server.Port)

	if listenErr != nil {
		logger.Fatalln(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	logger.Fatalln(server.Serve(listener))

}
