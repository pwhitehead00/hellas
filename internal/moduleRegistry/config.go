package moduleregistry

type server struct {
	CertSecretName string `yaml:"certSecretName"`
}

type GithubConfig struct {
	TokenSecretName    string `yaml:"tokenSecretName"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
	Protocol           string `yaml:"protocol"`
}

type Config struct {
	Server   server         `yaml:"server"`
	Registry map[string]any `yaml:"registry"`
}
