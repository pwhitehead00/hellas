package moduleregistry

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/google/go-github/v40/github"
	"github.com/ironhalo/hellas/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGitHubDownload(t *testing.T) {
	t.Run("Github Download Source with Prefix", func(t *testing.T) {
		mr := models.ModuleRegistry{
			InsecureSkipVerify: false,
			Protocol:           "https",
			Prefix:             "prefix",
		}

		expected := "git::https://github.com/my-namespace/prefix-happycloud-module?ref=v3.11.0"

		c := NewGitHubClient(mr)
		actual := c.Download("my-namespace", "module", "happycloud", "3.11.0")
		assert.Equal(t, expected, actual, "Validate github download source")
	})
	t.Run("Github Download Source without Prefix", func(t *testing.T) {
		mr := models.ModuleRegistry{
			InsecureSkipVerify: false,
			Protocol:           "https",
			Prefix:             "",
		}

		expected := "git::https://github.com/my-namespace/happycloud-module?ref=v3.11.0"

		c := NewGitHubClient(mr)
		actual := c.Download("my-namespace", "module", "happycloud", "3.11.0")
		assert.Equal(t, expected, actual, "Validate github download source")
	})
}

func TestGitHubClient(t *testing.T) {
	t.Run("Github Client Insecure TLS", func(t *testing.T) {
		mr := models.ModuleRegistry{
			InsecureSkipVerify: true,
			Protocol:           "https",
			Prefix:             "prefix",
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c := &http.Client{Transport: tr}
		expected := &GitHubClient{
			Client: github.NewClient(c),
			Config: &models.ModuleRegistry{
				InsecureSkipVerify: true,
				Protocol:           "https",
				Prefix:             "prefix",
			},
		}

		actual := NewGitHubClient(mr)
		assert.Equal(t, expected, actual, "Github clients should be equal")
	})
	t.Run("Github Client Secure TLS", func(t *testing.T) {
		mr := models.ModuleRegistry{
			InsecureSkipVerify: false,
			Protocol:           "https",
			Prefix:             "prefix",
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		}
		c := &http.Client{Transport: tr}
		expected := &GitHubClient{
			Client: github.NewClient(c),
			Config: &models.ModuleRegistry{
				InsecureSkipVerify: false,
				Protocol:           "https",
				Prefix:             "prefix",
			},
		}

		actual := NewGitHubClient(mr)
		assert.Equal(t, expected, actual, "Github clients should be equal")
	})
}

func TestGitHubVersions(t *testing.T) {
	t.Run("Github Versions", func(t *testing.T) {
		mr := models.ModuleRegistry{
			InsecureSkipVerify: true,
			Protocol:           "https",
			Prefix:             "prefix",
		}

		expected := models.ModuleVersions{
			Modules: []*models.ModuleProviderVersions{
				{
					Source: "my-namespace/prefix-provider-name",
					Versions: []*models.ModuleVersion{
						{
							Version: "1.0.0",
						},
						{
							Version: "1.0.1",
						},
					},
				},
			},
		}

		c := NewGitHubClient(mr)
		actual := c.Versions("my-namespace", "name", "provider", []string{"1.0.0", "1.0.1"})
		assert.Equal(t, expected, actual, "GitHub Versions should be the same")
	})
}
