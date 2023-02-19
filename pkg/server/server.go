package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/RedLabsPlatform/kube-shield/pkg/engine"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewServer(addr, certPath, keyPath string, ipv4 bool, engine *engine.Engine) *Server {

	logger := logrus.WithField("component", "server")

	return &Server{
		Address:   addr,
		ForceIPV4: ipv4,
		TLS: &TLSConfig{
			CertPath: certPath,
			KeyPath:  keyPath,
		},
		Engine: engine,
		Logger: logger,
	}
}

// Start HTTP server
func (s *Server) Start() error {

	var listener net.Listener
	listener, err := net.Listen("tcp6", s.Address)
	if err != nil {
		return err
	}

	if s.ForceIPV4 {
		listener, err = net.Listen("tcp4", s.Address)
	}

	http.HandleFunc("/validate", s.ServeValidate)

	s.Logger.Infoln("starting server on address", s.Address)
	err = http.ServeTLS(
		listener,
		nil,
		s.TLS.CertPath,
		s.TLS.KeyPath,
	)
	return err
}

// ServeValidate validates an admission request and then writes an admission review response to `w`
func (s *Server) ServeValidate(w http.ResponseWriter, r *http.Request) {

	var admissionResponse admissionv1.AdmissionResponse

	w.Header().Set("Content-Type", "application/json")

	admissionReview, err := getAdmissionReview(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.Engine.Run(admissionReview.Request)
	if err != nil {
		admissionResponse.Allowed = false
		admissionResponse.Result = &metav1.Status{
			Message: fmt.Sprintf("%v", err),
			Status:  fmt.Sprintf("%v", err),
		}
	} else {
		admissionResponse.Allowed = true
	}

	admissionResponse.UID = admissionReview.Request.UID
	admissionReview.Response = &admissionResponse

	json.NewEncoder(w).Encode(admissionReview)
}

// getAdmissionReview extracts an AdmissionReview from an http.Request if possible
func getAdmissionReview(r *http.Request) (*admissionv1.AdmissionReview, error) {

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
