package server

import (
	"load_balancer/config"
	"log/slog"
	"net/http"
)

func logInitialization(cfg config.ConfigLoadBalancer) {
	slog.Info("initializing load balancer",
		"load_balancer_name", cfg.Name,
		"listening_address", cfg.ListeningServer.Address,
		"listening_port", cfg.ListeningServer.Port,
		"algorithm", cfg.Algorithm,
		"backend_servers_length", len(cfg.Servers),
	)
}

func logSelectedServer(requestID string, serverIndex int, server Server) {
	slog.Info("selected server for request",
		"request_id", requestID,
		"server_index", serverIndex,
		"server_address", server.Address(),
	)
}

func logServerNotAlive(requestID string, serverIndex int, server Server) {
	slog.Warn("server is not alive, trying next one",
		"request_id", requestID,
		"server_index", serverIndex,
		"server_address", server.Address(),
	)
}

func logSuccessfullySelectedServer(requestID string, serverIndex int, server Server) {
	slog.Info("available server selected",
		"request_id", requestID,
		"server_index", serverIndex,
		"server_address", server.Address(),
	)
}

func logForwardingRequest(availableServer Server, req *http.Request) {
	slog.Info("forwarding request to server",
		"server_name", availableServer.Name(),
		"server_address", availableServer.Address(),
		"request_method", req.Method,
	)
}

func logReceivedResponse(req *http.Request, wrapper *ResponseWriterWrapper) {
	statusCode := wrapper.StatusCode()

	logLevel := slog.Info
	message := "request processing succeeded"

	if !(statusCode >= 200 && statusCode < 400) {
		logLevel = slog.Warn
		message = "request processing failed"
	}

	logLevel("response received from server",
		"status_code", statusCode,
		"request_method", req.Method,
		"request_status", message,
	)
}
