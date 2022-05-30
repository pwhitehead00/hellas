package moduleregistry

import (
	"fmt"

	"github.com/ironhalo/hellas/internal/models"
)

func Versions(namespace, name, provider, repo string, version []string) models.ModuleVersions {
	var m models.ModuleVersions
	var mv []*models.ModuleVersion

	for _, t := range version {
		o := models.ModuleVersion{
			Version: t,
		}
		mv = append(mv, &o)
	}

	mpv := models.ModuleProviderVersions{
		Source:   fmt.Sprintf("%s/%s", namespace, repo),
		Versions: mv,
	}

	m.Modules = append(m.Modules, &mpv)

	return m
}
