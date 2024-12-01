package moduleregistry

type wellKnown struct {
	Modules string `json:"modules.v1"`
}

type moduleVersions struct {
	Modules []*moduleProviderVersions `json:"modules"`
}

type moduleProviderVersions struct {
	Source   string           `json:"source"`
	Versions []*moduleVersion `json:"versions"`
}

type moduleVersion struct {
	Version string `json:"version"`
}

func newModuleVersions() moduleVersions {
	return moduleVersions{
		Modules: []*moduleProviderVersions{
			{},
		},
	}
}

func (mvs *moduleVersions) setSource(source string) {
	mvs.Modules[0].Source = source
}

func (mvs *moduleVersions) addVersion(version *string) {
	v := &moduleVersion{
		Version: *version,
	}

	mvs.Modules[0].Versions = append(mvs.Modules[0].Versions, v)
}
