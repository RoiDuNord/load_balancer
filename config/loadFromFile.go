package config

import (
	"errors"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

func loadFromFile(fileName string) (*Config, error) {
	slog.Debug("loading config file", "file", fileName)

	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn("config file not found", "file", fileName)
			return &Config{}, errors.New("config file not found")
		}
		slog.Error("error reading config file", "file", fileName, "error", err)
		return &Config{}, err
	}

	var cfg *Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		slog.Error("error unmarshaling config", "file", fileName, "error", err)
		return &Config{}, err
	}

	if err := cfg.IsEmpty(); err != nil {
		return &Config{}, err
	}

	slog.Info("config file loaded successfully", "file", fileName)
	return cfg, nil
}

func (cfg *Config) IsEmpty() error {
	if cfg == nil || cfg.ConfigLoadBalancer.Name == "" {
		slog.Warn("empty or invalid config")
		return errors.New("empty or invalid config error")
	}
	return nil
}

