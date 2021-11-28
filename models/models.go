package models

type WellKnown struct {
	Modules string `json:"modules.v1"`
}

type ModuleVersion struct {
	Modules []Versions `json:"modules"`
}

type Versions struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	Version string `json:"version"`
}
