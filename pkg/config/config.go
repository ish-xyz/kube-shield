package config

import (
	validator "github.com/go-playground/validator/v10"
)

type Config struct {
	Address        string `validate:"required"`
	Path           string `validate:"required"`
	TLSKey         string `validate:"required"`
	TLSCert        string `validate:"required"`
	Policies       string `validate:"required"`
	Register       bool
	Debug          bool
	MetricsAddress string `validate:"required"`
	MetricsPath    string `validate:"required"`
}

func NewConfig(policies, addr, path, key, cert string,
	register, debug bool,
	metricsAddr, metricsPath string) *Config {

	return &Config{
		Address:        addr,
		Path:           path,
		TLSKey:         key,
		TLSCert:        cert,
		Policies:       policies,
		Register:       register,
		Debug:          debug,
		MetricsAddress: metricsAddr,
		MetricsPath:    metricsPath,
	}
}

func (cfg *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(cfg)
}
