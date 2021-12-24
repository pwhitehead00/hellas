package moduleregistry

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/google/go-github/v40/github"
	"github.com/ironhalo/hellas/models"
	"github.com/stretchr/testify/assert"
)

func TestGitHubDownload(t *testing.T) {
	t.Run("GitHub downoad", func(t *testing.T) {
		expected := "git::https://github.com/my-namespace/terraform-happycloud-module?ref=v3.11.0"

		c := NewGitHubClient()
		got := c.Download("my-namespace", "module", "happycloud", "3.11.0")
		if got != expected {
			t.Fatalf("Expected %s, got %s", expected, got)
		}
	})
}

func TestGitHubClient(t *testing.T) {
	t.Run("Default github client", func(t *testing.T) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		}
		c := &http.Client{Transport: tr}
		expected := &GitHubClient{
			Client: github.NewClient(c),
		}

		actual := NewGitHubClient()
		assert.Equal(t, expected, actual, "GitHub Clients should be equal")
	})
}

func TestGitHubVersions(t *testing.T) {
	t.Run("GitHub versions", func(t *testing.T) {
		expected := models.ModuleVersions{
			Modules: []*models.ModuleProviderVersions{
				{
					Source: "my-namespace/terraform-provider-name",
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

		c := NewGitHubClient()
		actual := c.Versions("my-namespace", "name", "provider", []string{"1.0.0", "1.0.1"})
		assert.Equal(t, expected, actual, "GitHub Versions should be the same")
	})
}
