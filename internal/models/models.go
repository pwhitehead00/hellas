package models

type WellKnown struct {
	Modules string `json:"modules.v1"`
}

type ModuleVersions struct {
	Modules []*ModuleProviderVersions `json:"modules"`
}

type ModuleProviderVersions struct {
	Source   string           `json:"source"`
	Versions []*ModuleVersion `json:"versions"`
}

type ModuleVersion struct {
	Version string `json:"version"`
}

type ModuleRegistry struct {
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
	Protocol           string `json:"protocol"`
	Prefix             string `json:"prefix"`
}

type Config struct {
	ModuleBackend  string          `json:"moduleBackend"`
	ModuleRegistry *ModuleRegistry `json:"moduleRegistry"`
}
