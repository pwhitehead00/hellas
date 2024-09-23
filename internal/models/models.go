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

func NewModuleVersions() ModuleVersions {
	return ModuleVersions{
		Modules: []*ModuleProviderVersions{
			{},
		},
	}
}

func (mvs *ModuleVersions) SetSource(source string) {
	mvs.Modules[0].Source = source
}

func (mvs *ModuleVersions) AddVersion(version *string) {
	v := &ModuleVersion{
		Version: *version,
	}

	mvs.Modules[0].Versions = append(mvs.Modules[0].Versions, v)
}
