package server

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/RedLabsPlatform/kube-shield/pkg/config"
)

type Server struct {
	Http        *http.Server
	WebhookPath string
	Debug       bool
}

func NewServer(cfg *config.Config) (*Server, error) {

	cert, err := tls.LoadX509KeyPair(cfg.TLSCert, cfg.TLSKey)
	if err != nil {
		return nil, err
	}

	return &Server{
		Http: &http.Server{
			Addr: cfg.Address,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
		WebhookPath: cfg.Path,
		Debug:       cfg.Debug,
	}, nil
}

func (s *Server) Run() error {

	fmt.Printf("+%v", s)
	return nil
}
