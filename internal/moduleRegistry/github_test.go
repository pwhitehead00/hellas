package moduleregistry

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-github/v44/github"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGitHubDownload(t *testing.T) {
	t.Run("Github Download", func(t *testing.T) {
		c := &gitHubConfig{
			Protocol:   "https",
			RepoPrefix: "prefix",
		}

		r := NewGitHubRegistry(c)
		actual := r.Download("my-namespace", "module", "happycloud", "3.11.0")

		expected := "git::https://github.com/my-namespace/prefix-happycloud-module?ref=v3.11.0"
		assert.Equal(t, expected, actual, "Validate github download source")
	})
}

func TestGitHubRegistry(t *testing.T) {
	t.Run("Authenticated GitHub Repository", func(t *testing.T) {
		os.Setenv("TOKEN", "mytoken")

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		}

		expected := &GitHubRegistry{
			Client: authenticatedGitHubClient("mytoken", tr),
			Config: &gitHubConfig{},
		}

		actual := NewGitHubRegistry(&gitHubConfig{})
		assert.Equal(t, expected, actual)
	})
}

func TestAuthenticatedGitHubClient(t *testing.T) {
	t.Run("Authenticated Github Client Secure TLS", func(t *testing.T) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		}
		actual := authenticatedGitHubClient("token", tr)

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "token"},
		)

		client := &http.Client{Transport: tr}
		ctx := context.TODO()
		ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
		tc := oauth2.NewClient(ctx, ts)
		expected := github.NewClient(tc)
		assert.Equal(t, expected, actual, "Github clients should be equal")
	})
}

func TestGitHubValidation(t *testing.T) {
	t.Run("Github: Invalid Protocol", func(t *testing.T) {
		c := &gitHubConfig{
			Protocol: "foo",
		}

		mr := NewGitHubRegistry(c)
		err := mr.validate()

		assert.Equal(t, errors.New("Invalid protocol: foo"), err)
	})
}

func TestGitHubPath(t *testing.T) {
	t.Run("Path: With Repo Prefix", func(t *testing.T) {
		c := &gitHubConfig{
			RepoPrefix: "prefix",
		}

		mr := NewGitHubRegistry(c)
		actual := mr.Path("happycloud", "module")

		assert.Equal(t, "prefix-happycloud-module", actual)
	})

	t.Run("Path: Without Repo Prefix", func(t *testing.T) {
		c := &gitHubConfig{}

		mr := NewGitHubRegistry(c)
		actual := mr.Path("happycloud", "module")

		assert.Equal(t, "happycloud-module", actual)
	})
}
