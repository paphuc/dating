package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"dating/internal/app/api"
	"dating/internal/app/config"
	envconfig "dating/internal/pkg/config/env"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/health"
)

//go:embed swagger-ui/* swagger-ui/dating-api.yaml
var staticFiles embed.FS

func main() {
	logger := glog.New()
	stage := flag.String("stage", "dev", "set working environment")
	configPath := flag.String("config", "configs", "set configs path, default as: 'configs'")

	// error message
	em := config.ErrorMessage{ConfigPath: *configPath}
	if err := em.Init(); err != nil {
		logger.Errorf("failed to load error messages, err: %v", err)
	}

	// configs
	conf, err := config.New(*configPath, *stage)
	if err != nil {
		logger.Errorf("failed to load config, err: %v", err)
	}
	// envconfig.Load(&conf)

	var mongoConf config.MongoDB
	envconfig.Load(&mongoConf)
	if mongoConf.Address != "" {
		conf.Database.Mongo.Address = mongoConf.Address
	}

	if mongoConf.Database != "" {
		conf.Database.Mongo.Database = mongoConf.Database
	}
	logger.Infof("initializing HTTP routing...")
	router, err := api.Init(conf, em, staticFiles)
	if err != nil {
		logger.Panicf("failed to init routing, err: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", conf.HTTPServer.Address, conf.HTTPServer.Port)
	port := os.Getenv("PORT")
	if port != "" {
		addr = fmt.Sprintf("%s:%s", conf.HTTPServer.Address, port)
	}
	httpServer := http.Server{
		Addr:              addr,
		Handler:           router,
		ReadTimeout:       conf.HTTPServer.ReadTimeout,
		WriteTimeout:      conf.HTTPServer.WriteTimeout,
		ReadHeaderTimeout: conf.HTTPServer.ReadHeaderTimeout,
	}

	logger.Infof("starting HTTP server...")
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panicf("http.ListenAndServe() error: %v", err)
		}
	}()

	// tell the world that we're ready
	health.Ready()
	logger.Infof("HTTP Server is listening at: %v", addr)

	// gracefully shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), conf.HTTPServer.ShutdownTimeout)
	defer cancel()
	logger.Infof("shutting down http server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("http server shutdown with error: %v", err)
	}
	// shutdown background services goes here
}
