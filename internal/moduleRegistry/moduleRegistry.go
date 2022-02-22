package moduleregistry

import (
	"github.com/ironhalo/hellas/internal/models"
)

type Registry interface {
	GetVersions(namespace, name, provider string) ([]string, error)
	Versions(namespace, name, provider string, versions []string) models.ModuleVersions
	Download(namespace, name, provider, version string) string
}

func NewModuleRegistry(registryType string, mr models.ModuleRegistry) (r Registry) {
	switch registryType {
	case "github":
		r = NewGitHubClient(mr)
	default:
		panic("unsupported registy type")
	}
	return
}
