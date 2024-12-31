package server

import (
	"crypto/tls"
	"net/http"
	"time"
)

func NewServer(mux *http.ServeMux, skipTLSVerify bool, cert, key string) (*http.Server, error) {
	s := &http.Server{
		Addr:         ":8443",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
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
