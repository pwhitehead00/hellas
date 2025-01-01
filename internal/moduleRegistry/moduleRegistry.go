package moduleregistry

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	discoveryPath string = "GET /.well-known/terraform.json"
)

var (
	noRegistriesEnabled error = errors.New("no registries enabled")
)

type registry interface {
	Versions() http.HandlerFunc
	Download() http.HandlerFunc
}

// TODO: support S3
func NewModuleRegistry(config Config) (*http.ServeMux, error) {
	enabled := false
	mux := http.NewServeMux()
	mux.HandleFunc(discoveryPath, discovery)
	mux.HandleFunc("/healthcheck", healthCheck)

	if config.Registries.Github.Enabled {
		if err := config.Registries.Github.validate(); err != nil {
			return nil, fmt.Errorf("invalid github config: %w", err)
		}

		enabled = true
		r, err := NewGitHubRegistry(config.Registries.Github)
		if err != nil {
			return nil, err
		}

		mux.Handle("GET /v1/modules/{group}/{project}/github/versions", http.TimeoutHandler(r.Versions(), 10*time.Second, "failed to get github version: timeout"))
		mux.Handle("GET /v1/modules/{group}/{project}/github/{version}/download", http.TimeoutHandler(r.Download(), 10*time.Second, "failed to get github download: timeout"))
	}

	if !enabled {
		return nil, noRegistriesEnabled
	}

	return mux, nil
}

// Discovery process
// See https://developer.hashicorp.com/terraform/internals/remote-service-discovery#discovery-process
func discovery(w http.ResponseWriter, r *http.Request) {
	wk := wellKnown{
		Modules: "/v1/modules/",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wk); err != nil {
		log.Printf("json encoding failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("ok"); err != nil {
		log.Printf("json encoding failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
