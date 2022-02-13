package moduleregistry

import (
	"github.com/ironhalo/hellas/internal/models"
)

type Registry interface {
	GetVersions(namespace, name, provider string) ([]string, error)
	Versions(namespace, name, provider string, versions []string) models.ModuleVersions
	Download(namespace, name, provider, version string) string
}

func NewModuleRegistry(registryType string) (r Registry) {
	switch registryType {
	case "github":
		r = NewGitHubClient("/config/config.json")
	default:
		panic("unsupported registy type")
	}
	return
}
