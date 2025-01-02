package server

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"time"

	"github.com/pwhitehead00/hellas/internal/logging"
)

func NewServer(mux *http.ServeMux, skipTLSVerify bool, cert, key string) (*http.Server, error) {
	log := slog.NewLogLogger(logging.Handler, slog.LevelError)

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
		ErrorLog: log,
	}

	return s, nil
}
