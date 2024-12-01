package moduleregistry

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const (
	discoveryPath string = "GET /.well-known/terraform.json"
)

type Registry interface {
	Versions(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
}

// Build a new Registry interface
// TODO: support S3
func NewModuleRegistry(config Config) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	mux.HandleFunc(discoveryPath, discovery)
	mux.HandleFunc("/healthcheck", healthCheck)

	for r, rc := range config.Registry {
		switch r {
		case "github":
			// TODO: set defaults
			// TODO: validate config
			c, ok := rc.(GithubConfig)
			if !ok {
				return nil, errors.New("invalid github config")
			}

			r, err := NewGitHubRegistry(c)
			if err != nil {
				return nil, err
			}

			mux.HandleFunc("GET /v1/modules/{group}/{project}/github/versions", r.Versions)
			mux.HandleFunc("GET /v1/modules/{group}/{project}/github/{version}/download", r.Download)
			return mux, nil
		}
	}

	return nil, errors.New("Unsupported registry type")
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
