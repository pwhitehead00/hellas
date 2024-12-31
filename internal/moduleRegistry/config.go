package moduleregistry

import (
	"errors"
	"fmt"
)

type protocol string

var (
	protocolHTTPS protocol = "https"
	protocolSSH   protocol = "ssh"
)

var (
	invalidProtocol error = errors.New("invalid protocol")
)

type GithubConfig struct {
	TokenSecretName    string   `yaml:"tokenSecretName,omitempty"`
	Protocol           protocol `yaml:"protocol"`
	Enabled            bool     `yaml:"enabled"`
	InsecureSkipVerify bool     `yaml:"insecureSkipVerify"`
}

type Registries struct {
	Github GithubConfig `yaml:"github"`
}

type Config struct {
	Registries Registries `yaml:"registries"`
}

func (gh GithubConfig) Validate() error {
	switch gh.Protocol {
	case protocolHTTPS, protocolSSH:
		return nil
	}

	return fmt.Errorf("validation failure: %w: %s", invalidProtocol, gh.Protocol)
}
