package server

import (
	"context"
	"load_balancer/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

type Server interface {
	Name() string
	Address() string
	Port() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

type LoadBalancer struct {
	Name            string
	ListeningServer config.ListeningServer
	RoundRobinCount int
	Servers         []Server
	Algorithm       string
	HTTPServer      *http.Server
	Ctx             context.Context
}

func NewLoadBalancer(cfg config.ConfigLoadBalancer, ctx context.Context) *LoadBalancer {
	logInitialization(cfg)

	backendServers := cfg.Servers
	servers := make([]Server, len(backendServers))
	for i, backendServer := range backendServers {
		servers[i] = backendServer
	}

	defer slog.Info("load balancer initialized")

	return &LoadBalancer{
		Name:            cfg.Name,
		ListeningServer: cfg.ListeningServer,
		RoundRobinCount: 0,
		Servers:         servers,
		Algorithm:       cfg.Algorithm,
		HTTPServer: &http.Server{
			Addr: cfg.ListeningServer.Address + ":" + cfg.ListeningServer.Port,
		},
		Ctx: ctx,
	}
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.Servers[lb.RoundRobinCount%len(lb.Servers)]
	requestID := generateRequestID()

	logSelectedServer(requestID, lb.RoundRobinCount%len(lb.Servers), server)

	for !server.IsAlive() {
		lb.RoundRobinCount++
		server = lb.Servers[lb.RoundRobinCount%len(lb.Servers)]
		logServerNotAlive(requestID, lb.RoundRobinCount%len(lb.Servers), server)
	}

	logSuccessfullySelectedServer(requestID, lb.RoundRobinCount%len(lb.Servers), server)

	lb.RoundRobinCount++
	return server
}

func generateRequestID() string {
	return uuid.New().String()
}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	availableServer := lb.getNextAvailableServer()
	logForwardingRequest(availableServer, req)

	wrapper := &ResponseWriterWrapper{ResponseWriter: rw}
	availableServer.Serve(wrapper, req)

	logReceivedResponse(req, wrapper)
}

func (rw *ResponseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriterWrapper) StatusCode() int {
	return rw.statusCode
}

func (lb *LoadBalancer) GracefulShutdown() error {
	server := lb.HTTPServer
	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case <-shutdownSignals:
		slog.Info("received shutdown signal")
	case <-lb.Ctx.Done():
		slog.Info("context deadline exceeded")
	}

	if err := server.Shutdown(lb.Ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		if closeErr := server.Close(); closeErr != nil {
			slog.Error("Forced shutdown failed", "error", err)
			return closeErr
		}
		slog.Info("server forced shutdown complete")
		return err
	}

	slog.Info("server graceful shutdown complete")
	return nil
}
