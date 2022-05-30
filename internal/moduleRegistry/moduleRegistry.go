package moduleregistry

import (
	"errors"
	"fmt"

	"github.com/ironhalo/hellas/internal/models"
)

type Registry interface {
	GetVersions(namespace, name, provider string) ([]string, error)
	Versions(namespace, name, provider string, versions []string) models.ModuleVersions
	Download(namespace, name, provider, version string) string
	Repo(provider, name string) string
	validate() error
}

func NewModuleRegistry(registryType *string, config []byte) (Registry, error) {
	var r Registry

	switch *registryType {
	case "github":
		c, err := newGitHubConfig(config)
		if err != nil {
			return nil, err
		}

		r, err = NewGitHubRegistry(c)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported registy type: %s", *registryType))
	}

	if err := r.validate(); err != nil {
		return nil, err
	}

	return r, nil
}
