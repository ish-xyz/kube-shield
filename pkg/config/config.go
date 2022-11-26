package config

import "fmt"

type Config struct {
	Policies        string
	RegisterWebhook bool
	Debug           bool
	MetricsAddress  string
}

func NewConfig(policies string, regWebhook, debug bool, metricsAddr string) *Config {
	return &Config{
		Policies:        policies,
		RegisterWebhook: regWebhook,
		Debug:           debug,
		MetricsAddress:  metricsAddr,
	}
}

func (cfg *Config) Validate() error {
	return fmt.Errorf("validation failed")
}
