package moduleregistry

import (
	"errors"
	"fmt"
	"log"

	"github.com/pwhitehead00/hellas/internal/models"
)

type Registry interface {
	// ListVersions(namespace, name, provider string) ([]string, error)
	ListVersions(group, path string) (*models.ModuleVersions, error)
	// Download(name, provider, version string) string
}

// Build a new Registry interface
func NewModuleRegistry(config Config) (Registry, error) {
	for r, rc := range config.Registry {
		switch r {
		case "github":
			log.Println("running github")
			c, ok := rc.(GithubConfig)
			if !ok {
				return nil, errors.New("invalid github config")
			}
			return NewGitHubRegistry(c)
		}
	}

	return nil, errors.New(fmt.Sprint("Unsupported registry type"))
}
