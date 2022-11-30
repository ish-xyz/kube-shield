package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RedLabsPlatform/kube-shield/pkg/config"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
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

	http.HandleFunc("/validate", ServeValidate)

	// Run tls server
	err := s.Http.ListenAndServeTLS("", "")
	return err
}

// ServeValidatePods validates an admission request and then writes an admission review to `w`
func ServeValidate(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("uri", r.RequestURI)
	logger.Debug("received validation request")

	payload, err := parseRequest(*r)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Warnln(payload)
}

// parseRequest extracts an AdmissionReview from an http.Request if possible
func parseRequest(r http.Request) (*admissionv1.AdmissionReview, error) {

	var a admissionv1.AdmissionReview

	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("Content-Type: %q should be %q",
			r.Header.Get("Content-Type"), "application/json")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read body request")
	}

	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("could not parse admission review request: %v", err)
	}

	if a.Request == nil {
		return nil, fmt.Errorf("admission review can't be used: Request field is nil")
	}

	return &a, nil
}
