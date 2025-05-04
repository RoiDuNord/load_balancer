package config

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	ConfigLoadBalancer ConfigLoadBalancer `yaml:"load_balancer"`
}

type ConfigLoadBalancer struct {
	Name            string                `yaml:"name"`
	ListeningServer ListeningServer       `yaml:"listen"`
	Servers         []*SimpleServerConfig `yaml:"backend_servers"`
	Algorithm       string                `yaml:"algorithm"`
}

type ListeningServer struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

type SimpleServerConfig struct {
	N     string `yaml:"name"`
	Addr  string `yaml:"address"`
	P     string `yaml:"port"`
	Proxy *httputil.ReverseProxy
}

func ParseConfig() (Config, error) {
	slog.Info("checking for config file", "file", "config.yaml")
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		slog.Error("config file not found", "file", "config.yaml")
		return Config{}, fmt.Errorf("config.yaml not found")
	}

	slog.Info("loading config file", "file", "config.yaml")
	cfg, err := loadFromFile("config.yaml")
	if err != nil {
		slog.Error("Error loading config file", "file", "config.yaml", "error", err)
		return Config{}, err
	}

	slog.Info("adding reverse proxy")
	if err = cfg.addReverseProxy(); err != nil {
		slog.Error("adding reverse proxy failed", "error", err)
		return Config{}, err
	}

	slog.Info("validating config load balancer")
	if err = cfg.ConfigLoadBalancer.validate(); err != nil {
		slog.Error("config load balancer validation failed", "error", err)
		return Config{}, err
	}

	slog.Info("config loaded and validated")
	return *cfg, nil
}

func (cfg *Config) addReverseProxy() error {
	for i, server := range cfg.ConfigLoadBalancer.Servers {
		slog.Info("parsing backend server URL", "index", i, "addr", server.Addr)
		serverUrl, err := url.Parse(server.Addr)
		if err != nil {
			slog.Error("error parsing backend server URL", "addr", server.Addr, "error", err)
			return err
		}
		slog.Info("creating reverse proxy for backend server", "addr", server.Addr)
		server.Proxy = httputil.NewSingleHostReverseProxy(serverUrl)
	}
	return nil
}

func (cfglb ConfigLoadBalancer) validate() error {
	if cfglb.Name == "" {
		slog.Warn("missing name for load balancer")
		return fmt.Errorf("missing name")
	}
	if cfglb.ListeningServer.Address == "" || cfglb.ListeningServer.Port == "" {
		slog.Warn("invalid listening server address or port")
		return fmt.Errorf("invalid listening server address or port")
	}
	if len(cfglb.Servers) == 0 {
		slog.Warn("no backend servers configured")
		return fmt.Errorf("no backend servers configured")
	}

	return nil
}

func (s *SimpleServerConfig) Name() string    { return s.N }
func (s *SimpleServerConfig) Address() string { return s.Addr }
func (s *SimpleServerConfig) Port() string    { return s.P }

func (s *SimpleServerConfig) IsAlive() bool {
	var flag bool
	slog.Info("checking server health", "addr", s.Addr)
	resp, err := http.Head(s.Addr)
	if err != nil {
		slog.Warn("failed to check backend server status", slog.String("addr", s.Addr), slog.String("error", err.Error()))
		return flag
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	if statusCode >= 200 && statusCode < 300 {
		slog.Info("backend server is alive", "addr", s.Addr)
		flag = true
	} else {
		slog.Warn("backend server is not alive", slog.String("addr", s.Addr), slog.String("status_code", strconv.Itoa(resp.StatusCode)))
	}

	return flag
}

func (s *SimpleServerConfig) Serve(rw http.ResponseWriter, req *http.Request) {
	s.Proxy.ServeHTTP(rw, req)
}
