package moduleregistry

import (
	"crypto/tls"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v44/github"
	"github.com/ironhalo/hellas/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGitHubDownload(t *testing.T) {
	t.Run("Github Download Source with Prefix", func(t *testing.T) {
		c := &gitHubConfig{
			Protocol: "https",
			Prefix:   "prefix",
		}

		r, _ := NewGitHubRegistry(c)
		actual := r.Download("my-namespace", "module", "happycloud", "3.11.0")

		expected := "git::https://github.com/my-namespace/prefix-happycloud-module?ref=v3.11.0"
		assert.Equal(t, expected, actual, "Validate github download source")
	})

	t.Run("Github Download Source without Prefix", func(t *testing.T) {
		c := &gitHubConfig{
			Protocol: "https",
		}

		r, _ := NewGitHubRegistry(c)
		actual := r.Download("my-namespace", "module", "happycloud", "3.11.0")

		expected := "git::https://github.com/my-namespace/happycloud-module?ref=v3.11.0"
		assert.Equal(t, expected, actual, "Validate github download source")
	})
}

func TestGitHubClient(t *testing.T) {
	t.Run("Github Client Insecure TLS", func(t *testing.T) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		expected := &GitHubRegistry{
			Client: github.NewClient(client),
			Config: &gitHubConfig{
				InsecureSkipVerify: true,
				Protocol:           "https",
				Prefix:             "prefix",
			},
		}

		c := &gitHubConfig{
			InsecureSkipVerify: true,
			Protocol:           "https",
			Prefix:             "prefix",
		}

		actual, _ := NewGitHubRegistry(c)
		assert.Equal(t, expected, actual, "Github clients should be equal")
	})

	t.Run("Github Client Secure TLS", func(t *testing.T) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		}
		client := &http.Client{Transport: tr}
		expected := &GitHubRegistry{
			Client: github.NewClient(client),
			Config: &gitHubConfig{
				InsecureSkipVerify: false,
				Protocol:           "https",
				Prefix:             "prefix",
			},
		}

		c := &gitHubConfig{
			InsecureSkipVerify: false,
			Protocol:           "https",
			Prefix:             "prefix",
		}

		actual, _ := NewGitHubRegistry(c)
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

		c := &gitHubConfig{
			Prefix: "prefix",
		}

		mr, _ := NewGitHubRegistry(c)
		actual := mr.Versions("my-namespace", "name", "provider", []string{"1.0.0", "1.0.1"})
		assert.Equal(t, expected, actual, "GitHub Versions should be the same")
	})
}

func TestGitHubValidation(t *testing.T) {
	t.Run("Github: Invalid Protocol", func(t *testing.T) {
		c := &gitHubConfig{
			Protocol: "foo",
		}

		mr, _ := NewGitHubRegistry(c)
		err := mr.validate()

		assert.Equal(t, errors.New("Invalid protocol: foo"), err)
	})
}

func TestGitHubRepo(t *testing.T) {
	t.Run("Repo: With Prefix", func(t *testing.T) {
		c := &gitHubConfig{
			Prefix: "prefix",
		}

		mr, _ := NewGitHubRegistry(c)
		actual := mr.Repo("happycloud", "module")

		assert.Equal(t, "prefix-happycloud-module", actual)
	})

	t.Run("Repo: Without Prefix", func(t *testing.T) {
		c := &gitHubConfig{}

		mr, _ := NewGitHubRegistry(c)
		actual := mr.Repo("happycloud", "module")

		assert.Equal(t, "happycloud-module", actual)
	})
}
