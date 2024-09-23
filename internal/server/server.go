package server

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"log"

	"github.com/pwhitehead00/hellas/internal/models"
)

const (
	discoveryPath string = "GET /.well-known/terraform.json"
)

func NewServer(mux *http.ServeMux, skipTLSVerify bool, cert, key string) (*http.Server, error) {
	s := &http.Server{
		Addr:         ":8443",
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: skipTLSVerify,
			GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
				caFiles, err := tls.LoadX509KeyPair(cert, key)
				if err != nil {
					return nil, err
				}

				return &caFiles, nil
			},
		},
	}

	return s, nil
}

func WellKnown(w http.ResponseWriter, r *http.Request) {
	wk := models.WellKnown{
		Modules: "/terraform/modules/v1/",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wk); err != nil {
		log.Printf("json encoding failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
