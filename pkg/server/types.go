package server

import (
	"github.com/RedLabsPlatform/kube-shield/pkg/engine"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Address   string
	ForceIPV4 bool
	TLS       *TLSConfig
	Engine    *engine.Engine
	Logger    *logrus.Entry
}

type TLSConfig struct {
	CertPath string
	KeyPath  string
}
