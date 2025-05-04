package main

import (
	"context"
	"fmt"
	"load_balancer/config"
	"load_balancer/logger"
	"load_balancer/server"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	logFile, err := logger.Init()
	if err != nil {
		return
	}
	defer logger.Close(logFile)

	cfg, err := config.ParseConfig()
	if err != nil {
		return
	}

	if err := runApp(cfg); err != nil {
		slog.Error("application encountered an error", "error", err)
		return
	}
}

func runApp(cfg config.Config) error {
	slog.Info("running App")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	deadline, ok := ctx.Deadline()
	if ok {
		remaining := time.Until(deadline)

		slog.Info("deadline set",
			"deadline_time", deadline.Format("15:04:05"),
			"remaining_time", fmt.Sprintf("%.0f seconds", remaining.Seconds()),
		)
	} else {
		slog.Warn("deadline not set")
	}
	defer cancel()

	lb := server.NewLoadBalancer(cfg.ConfigLoadBalancer, ctx)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.ServeProxy(rw, req)
	}

	http.HandleFunc("/", handleRedirect)

	go func() {
		if err := startHTTPServer(lb.HTTPServer, lb.ListeningServer.Port); err != nil {
			slog.Error("HTTP server exited with error", "error", err)
		}
	}()

	return lb.GracefulShutdown()
}

func startHTTPServer(server *http.Server, port string) error {
	slog.Info(fmt.Sprintf("starting HTTP server on port: %s", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("error starting HTTP server", "error", err)
		return err
	}
	return nil
}
