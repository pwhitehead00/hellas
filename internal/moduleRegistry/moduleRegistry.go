package moduleregistry

import (
	"github.com/ironhalo/hellas/models"
)

type Registry interface {
	Versions(namespace, name, provider string) models.ModuleVersions
	Download(namespace, name, provider, version string) string
}

func NewModuleRegistry(registryType string) (r Registry) {
	switch registryType {
	case "github":
		r = NewGitHubClient()
	default:
		panic("unsupported registy type")
	}
	return
}
