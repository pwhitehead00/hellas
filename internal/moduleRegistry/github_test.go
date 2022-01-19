package moduleregistry

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-github/v40/github"
	"github.com/ironhalo/hellas/models"
	"github.com/stretchr/testify/assert"
)

func TestGitHubDownload(t *testing.T) {
	t.Run("Github Download Source with Prefix", func(t *testing.T) {
		expected := "git::https://github.com/my-namespace/prefix-happycloud-module?ref=v3.11.0"

		os.Setenv("CONFIG", "../../test/github-default.json")
		c := NewGitHubClient()
		actual := c.Download("my-namespace", "module", "happycloud", "3.11.0")
		assert.Equal(t, expected, actual, "Validate github download source")
	})
	// t.Run("Github Download Source without Prefix", func(t *testing.T) {
	// 	expected := "git::https://github.com/my-namespace/happycloud-module?ref=v3.11.0"

	// 	os.Setenv("CONFIG", "../../test/github-noprefix.json")
	// 	c := NewGitHubClient()
	// 	actual := c.Download("my-namespace", "module", "happycloud", "3.11.0")
	// 	assert.Equal(t, expected, actual, "Validate github download source")
	// })
}

func TestGitHubClient(t *testing.T) {
	t.Run("Default Github Client", func(t *testing.T) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c := &http.Client{Transport: tr}
		expected := &GitHubClient{
			Client: github.NewClient(c),
			Config: &GitHubConfig{
				InsecureSkipVerify: true,
				Protocol:           "https",
				Prefix:             "prefix",
			},
		}

		os.Setenv("CONFIG", "../../test/github-default.json")
		actual := NewGitHubClient()
		assert.Equal(t, expected, actual, "Github clients should be equal")
	})
}

func TestGitHubVersions(t *testing.T) {
	t.Run("Github Versions", func(t *testing.T) {
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

		os.Setenv("CONFIG", "../../test/github-default.json")
		c := NewGitHubClient()
		actual := c.Versions("my-namespace", "name", "provider", []string{"1.0.0", "1.0.1"})
		assert.Equal(t, expected, actual, "GitHub Versions should be the same")
	})
}
