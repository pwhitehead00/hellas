package moduleregistry

import (
	"crypto/tls"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v44/github"
	"github.com/stretchr/testify/assert"
)

func TestGitHubDownload(t *testing.T) {
	t.Run("Github Download", func(t *testing.T) {
		c := &gitHubConfig{
			Protocol:   "https",
			RepoPrefix: "prefix",
		}

		r, _ := NewGitHubRegistry(c)
		actual := r.Download("my-namespace", "module", "happycloud", "3.11.0")

		expected := "git::https://github.com/my-namespace/prefix-happycloud-module?ref=v3.11.0"
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
				RepoPrefix:         "repoPrefix",
			},
		}

		c := &gitHubConfig{
			InsecureSkipVerify: true,
			Protocol:           "https",
			RepoPrefix:         "repoPrefix",
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
				RepoPrefix:         "repoPrefix",
			},
		}

		c := &gitHubConfig{
			InsecureSkipVerify: false,
			Protocol:           "https",
			RepoPrefix:         "repoPrefix",
		}

		actual, _ := NewGitHubRegistry(c)
		assert.Equal(t, expected, actual, "Github clients should be equal")
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

func TestGitHubPath(t *testing.T) {
	t.Run("Path: With Repo Prefix", func(t *testing.T) {
		c := &gitHubConfig{
			RepoPrefix: "prefix",
		}

		mr, _ := NewGitHubRegistry(c)
		actual := mr.Path("happycloud", "module")

		assert.Equal(t, "prefix-happycloud-module", actual)
	})

	t.Run("Path: Without Repo Prefix", func(t *testing.T) {
		c := &gitHubConfig{}

		mr, _ := NewGitHubRegistry(c)
		actual := mr.Path("happycloud", "module")

		assert.Equal(t, "happycloud-module", actual)
	})
}
