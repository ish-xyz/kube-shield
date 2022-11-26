package server

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
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

	http.HandleFunc("/validate-pods", ServeValidatePods)

	// run tls server
	err := s.Http.ListenAndServeTLS("", "")
	return err
}

// ServeValidatePods validates an admission request and then writes an admission
// review to `w`
func ServeValidatePods(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("uri", r.RequestURI)
	logger.Debug("received validation request")

	in, err := parseRequest(*r)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Debug(in)

	// write validation logic

	// adm := admission.Admitter{
	// 	Logger:  logger,
	// 	Request: in.Request,
	// }

	// out, err := adm.ValidatePodReview()
	// if err != nil {
	// 	e := fmt.Sprintf("could not generate admission response: %v", err)
	// 	logger.Error(e)
	// 	http.Error(w, e, http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// jout, err := json.Marshal(out)
	// if err != nil {
	// 	e := fmt.Sprintf("could not parse admission response: %v", err)
	// 	logger.Error(e)
	// 	http.Error(w, e, http.StatusInternalServerError)
	// 	return
	// }

	// logger.Debug("sending response")
	// logger.Debugf("%s", jout)
	// fmt.Fprintf(w, "%s", jout)
}

// parseRequest extracts an AdmissionReview from an http.Request if possible
func parseRequest(r http.Request) (*admissionv1.AdmissionReview, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("Content-Type: %q should be %q",
			r.Header.Get("Content-Type"), "application/json")
	}

	bodybuf := new(bytes.Buffer)
	bodybuf.ReadFrom(r.Body)
	body := bodybuf.Bytes()

	if len(body) == 0 {
		return nil, fmt.Errorf("admission request body is empty")
	}

	var a admissionv1.AdmissionReview

	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("could not parse admission review request: %v", err)
	}

	if a.Request == nil {
		return nil, fmt.Errorf("admission review can't be used: Request field is nil")
	}

	return &a, nil
}
