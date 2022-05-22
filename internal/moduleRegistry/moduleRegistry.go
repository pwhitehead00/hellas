package moduleregistry

import (
	"errors"
	"fmt"
)

type Registry interface {
	ListVersions(namespace, name, provider string) ([]string, error)
	Download(namespace, name, provider, version string) string
	Path(provider, name string) string
	validate() error
}

// Build a new Registry interface
func NewModuleRegistry(registryType *string, config []byte) (Registry, error) {
	var r Registry

	switch *registryType {
	case "github":
		c, err := newGitHubConfig(config)
		if err != nil {
			return nil, err
		}

		r = NewGitHubRegistry(c)
	case "gitlab":
		c, err := newGitLabConfig(config)
		if err != nil {
			return nil, err
		}

		r = NewGitLabRegistry(c)
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported registy type: %s", *registryType))
	}

	if err := r.validate(); err != nil {
		return nil, err
	}

	return r, nil
}
