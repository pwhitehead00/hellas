package models

type WellKnown struct {
	Modules string `json:"modules.v1"`
}

type Modules struct {
	Modules []struct {
		Versions []*map[string]string `json:"versions"`
	} `json:"modules"`
}

type Versions struct {
	Versions []*map[string]string `json:"versions"`
}
