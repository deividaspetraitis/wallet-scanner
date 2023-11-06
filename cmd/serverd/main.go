package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deividaspetraitis/wallet-screener/config"
	"github.com/deividaspetraitis/wallet-screener/errors"
	ihttp "github.com/deividaspetraitis/wallet-screener/http"
	"github.com/deividaspetraitis/wallet-screener/log"
	"github.com/deividaspetraitis/wallet-screener/riskprovider"

	immudb "github.com/codenotary/immudb/pkg/client"
)

var shutdowntimeout = time.Duration(5) * time.Second

// program flags
var (
	cfgPath string
)

// initialise program state
func init() {
	flag.StringVar(&cfgPath, "config", os.Getenv("config"), "PATH to .env configuration file")
}

// main program entry point.
func main() {
	flag.Parse()

	logger := log.Default()

	cfg, err := config.New(cfgPath)
	if err != nil {
		logger.WithError(err).Error("parsing configuration file")
	}

	if err := run(cfg, logger); err != nil {
		logger.WithError(err).Error("unable to start service")
		os.Exit(1)
	}
}

func run(cfg *config.Config, logger log.Logger) error {
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Construct services

	// even though the server address and port are defaults, setting them as a reference
	opts := immudb.DefaultOptions().WithAddress(cfg.Database.Host).WithPort(cfg.Database.Port)

	// construct a new immudb client
	immudbclient := immudb.NewClient().WithOptions(opts)

	// connect with immudb server (user, password, database)
	err := immudbclient.OpenSession(context.Background(), []byte(cfg.Database.Username), []byte(cfg.Database.Password), cfg.Database.Database)
	if err != nil {
		return errors.Wrap(err, "unable connect to immudb instance")
	}

	// Construct risk provider API client.
	httpclient, err := ihttp.NewClient(
		"https://api.blockmate.io/v1",                  // API URL ( unlikely to change )
		ihttp.WithHeader("Accept", "application/json"), // speaks with JSON
	)
	if err != nil {
		logger.WithError(err).Fatal("unable to construct Blockmate risk provider client")
	}

	// Construct specific risk provider.
	riskprovider, err := riskprovider.NewBlockMate(cfg.RiskProvider.Blockmate.APIKey, httpclient)
	if err != nil {
		logger.WithError(err).Fatal("unable to construct Blockmate risk provider")
	}

	// =========================================================================
	// Start HTTP server

	api := http.Server{
		Addr:    cfg.HTTP.Address,
		Handler: ihttp.API(shutdown, cfg.HTTP, logger, riskprovider, immudbclient),
	}

	go func() {
		logger.Printf("http server listening on %s", cfg.HTTP.Address)
		serverErrors <- api.ListenAndServe()
	}()

	// ========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		logger.Printf("http server start shutdown caused by %v", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), shutdowntimeout)
		defer cancel()

		if err := immudbclient.CloseSession(ctx); err != nil {
			logger.WithError(err).Error("graceful shutdown did not complete")
		}

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			logger.WithError(err).Error("graceful shutdown did not complete")
			api.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
