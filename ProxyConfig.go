package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ProxyConfig struct {
	Port            int           `json:"port"`
	AdminPort       int           `json:"admin_port"`
	Strategy        string        `json:"strategy"`
	HealthCheckFreq time.Duration `json:"health_check_frequency"`
	Backends        []string      `json:"backends"`
}

type proxyConfigJSON struct {
	Port            int      `json:"port"`
	AdminPort       int      `json:"admin_port"`
	Strategy        string   `json:"strategy"`
	HealthCheckFreq string   `json:"health_check_frequency"`
	Backends        []string `json:"backends"`
}

func LoadConfig(filename string) (*ProxyConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var configJSON proxyConfigJSON
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&configJSON); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	config := &ProxyConfig{
		Port:     configJSON.Port,
		AdminPort: configJSON.AdminPort,
		Strategy: configJSON.Strategy,
		Backends: configJSON.Backends,
	}

	if configJSON.HealthCheckFreq != "" {
		duration, err := time.ParseDuration(configJSON.HealthCheckFreq)
		if err != nil {
			return nil, fmt.Errorf("invalid health_check_frequency format: %w", err)
		}
		config.HealthCheckFreq = duration
	}

	if config.Port == 0 {
		config.Port = 9000
	}
	if config.AdminPort == 0 {
		config.AdminPort = 8081
	}
	if config.HealthCheckFreq == 0 {
		config.HealthCheckFreq = 30 * time.Second
	}
	if config.Strategy == "" {
		config.Strategy = "round-robin"
	}

	return config, nil
}
