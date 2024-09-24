package moduleregistry

import (
	"errors"
	"net/http"
)

type Registry interface {
	Versions(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
}

// Build a new Registry interface
// TODO: support S3
func NewModuleRegistry(config Config) (Registry, error) {
	for r, rc := range config.Registry {
		switch r {
		case "github":
			// TODO: set defaults
			// TODO: validate config
			c, ok := rc.(GithubConfig)
			if !ok {
				return nil, errors.New("invalid github config")
			}
			return NewGitHubRegistry(c)
		}
	}

	return nil, errors.New("Unsupported registry type")
}
