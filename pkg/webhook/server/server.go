package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/RedLabsPlatform/kube-shield/pkg/webhook/engine"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
)

type Server struct {
	Engine *engine.Engine
}

func (s *Server) Start() {

	http.HandleFunc("/validate", s.ServeValidate)

	// Run tls server
	addr := ":8000" //TODO: should be a parameter
	logrus.Infoln("serving requests on" + addr)
	// TODO: temporary for testing with minikube (it doesn't support tunneling over IPV6)
	// logrus.Fatal(http.ListenAndServeTLS("0.0.0.0:8001", "./certs/server.crt", "./certs/server.key", nil))
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	logrus.Fatal(http.ServeTLS(l, nil, "./certs/server.crt", "./certs/server.key"))

}

// ServeValidatePods validates an admission request and then writes an admission review to `w`
func (s *Server) ServeValidate(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("uri", r.RequestURI)
	logger.Info("received validation request")

	payload, err := getAdmissionReview(r)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.Engine.Run(payload)
	w.WriteHeader(500)
	fmt.Fprint(w, "internal server error")
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
