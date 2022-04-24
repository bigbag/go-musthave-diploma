package main

import (
	"io/ioutil"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-diploma/internal/app"
	"github.com/bigbag/go-musthave-diploma/internal/config"
)

func getLogger(cfg *config.Config, baseLogger *stdLog.Logger) logrus.StdLogger {
	logger := logrus.New()

	if level, err := logrus.ParseLevel(cfg.Logger.Level); err == nil {
		baseLogger.Printf("Logger level set to %s\n", level.String())
		logger.SetLevel(level)
	} else {
		baseLogger.Printf("Failed to parse log level: %s\n", err.Error())
		logger.SetLevel(logrus.ErrorLevel)
	}

	switch cfg.Logger.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		logger.SetOutput(ioutil.Discard)
	}

	switch cfg.Logger.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	return logger
}

// shutdown implements all graceful shutdown logic
func shutdown(_ os.Signal, server *app.Server, l logrus.StdLogger) {
	l.Println("Shutting down...")
	if err := server.Stop(); err != nil {
		l.Printf("Error stopping api: %v", err)
	}
	l.Println("Running cleanup tasks...")
	l.Println("Fiber was successful shutdown.")
	os.Exit(0)
}

func main() {
	baseLogger := stdLog.New(os.Stdout, "", 0)

	cfg, err := config.New()
	if err != nil {
		baseLogger.Fatalf("Failed to initialize config: %v\n", err)
	}

	var (
		l      = getLogger(cfg, baseLogger)
		server = app.New(l.(logrus.FieldLogger), cfg)
	)

	// start HTTP API server
	go func() {
		l.Printf("API server listening on %s", cfg.Server.Listen)
		if err := server.Start(cfg.Server.Listen); err != nil && err != http.ErrServerClosed {
			l.Fatalf("Failed to start api server: %v", err)
		}
	}()

	// listen for exit signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	shutdown(<-c, server, l)
}
