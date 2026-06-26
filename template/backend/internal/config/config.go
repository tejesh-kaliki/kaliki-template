// Package config loads runtime configuration from a YAML file with env overrides.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server        ServerConfig        `yaml:"server"`
	Database      DatabaseConfig      `yaml:"database"`
	Observability ObservabilityConfig `yaml:"observability"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	URL string `yaml:"url"`
}

// ObservabilityConfig is wired in by default. Leaving Endpoint empty selects a
// no-op tracer, so the app runs with zero external dependencies until you point
// it at a real OTel collector.
type ObservabilityConfig struct {
	ServiceName string `yaml:"service_name"`
	Endpoint    string `yaml:"endpoint"`
}

// Load reads the YAML at path, then applies env overrides for anything that is
// commonly injected by the deployment environment.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Server:        ServerConfig{Port: "8080"},
		Observability: ObservabilityConfig{ServiceName: "backend"},
	}

	if b, err := os.ReadFile(path); err == nil {
		if err := yaml.Unmarshal(b, cfg); err != nil {
			return nil, err
		}
	}

	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.Database.URL = v
	}
	if v := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); v != "" {
		cfg.Observability.Endpoint = v
	}

	return cfg, nil
}
