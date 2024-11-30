package moduleregistry

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/pwhitehead00/hellas/internal/models"
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
	mux.HandleFunc(discoveryPath, wellKnown)

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

			mux.HandleFunc("GET /v1/modules/github/{group}/{project}/versions", r.Versions)
			mux.HandleFunc("GET /v1/modules/github/{group}/{project}/{version}/download", r.Download)
			return mux, nil
		}
	}

	return nil, errors.New("Unsupported registry type")
}

func wellKnown(w http.ResponseWriter, r *http.Request) {
	wk := models.WellKnown{
		Modules: "/terraform/modules/v1/",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(wk); err != nil {
		log.Printf("json encoding failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
