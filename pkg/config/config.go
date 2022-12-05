package config

import (
	validator "github.com/go-playground/validator/v10"
)

type Config struct {
	Address        string `validate:"required"`
	Path           string `validate:"required"`
	TLSKey         string `validate:"required"`
	TLSCert        string `validate:"required"`
	MetricsAddress string `validate:"required"`
	MetricsPath    string `validate:"required"`
	Debug          bool
}

func NewConfig(addr, path, key, cert string, debug bool, metricsAddr, metricsPath string) *Config {

	return &Config{
		Address:        addr,
		Path:           path,
		TLSKey:         key,
		TLSCert:        cert,
		MetricsAddress: metricsAddr,
		MetricsPath:    metricsPath,
		Debug:          debug,
	}
}

func (cfg *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(cfg)
}
