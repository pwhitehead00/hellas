package moduleregistry

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

type GitHubRegistry struct {
	Client *github.Client
	Config *gitHubConfig
}

type gitHubConfig struct {
	// Accept self signed certs
	InsecureSkipVerify bool `json:"insecureSkipVerify"`

	// The protocol used when checking out terraform modules.  Can be either
	// https or ssh
	Protocol string `json:"protocol"`

	// A string prefix
	RepoPrefix string `json:"repoPrefix"`
}

func newGitHubConfig(file []byte) (*gitHubConfig, error) {
	var config gitHubConfig

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Bulid an unauthenticated GitHub client
func unauthenticatedGitHubClient(tr *http.Transport) *github.Client {
	client := &http.Client{Transport: tr}
	gitHubClient := github.NewClient(client)

	return gitHubClient
}

// Bulid an authenticated GitHub client
func authenticatedGitHubClient(token string, tr *http.Transport) *github.Client {
	client := &http.Client{Transport: tr}
	ctx := context.TODO()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, client)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	gitHubClient := github.NewClient(tc)
	return gitHubClient
}

// New GitHub module registry
func NewGitHubRegistry(config *gitHubConfig) Registry {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify},
	}

	token, ok := os.LookupEnv("TOKEN")

	if ok {
		client := authenticatedGitHubClient(token, tr)
		return &GitHubRegistry{
			Client: client,
			Config: config,
		}
	}

	client := unauthenticatedGitHubClient(tr)
	return &GitHubRegistry{
		Client: client,
		Config: config,
	}
}

// Helper function to build the GitHub repo path
func (gh *GitHubRegistry) Path(provider, name string) string {
	if gh.Config.RepoPrefix == "" {
		return fmt.Sprintf("%s-%s", provider, name)
	}

	return fmt.Sprintf("%s-%s-%s", gh.Config.RepoPrefix, provider, name)
}

// List all tags for a GitHub registry
func (gh *GitHubRegistry) ListVersions(namespace, name, provider string) ([]string, error) {
	var allTags []*github.RepositoryTag
	var versions []string

	opt := &github.ListOptions{
		PerPage: 100,
	}

	repo := gh.Path(provider, name)

	for {
		tags, resp, err := gh.Client.Repositories.ListTags(context.Background(), namespace, repo, opt)
		if err != nil {
			return nil, err
		}

		allTags = append(allTags, tags...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, v := range allTags {
		versions = append(versions, *v.Name)
	}
	return versions, nil
}

// Download source code for a specific module version
// See https://www.terraform.io/internals/module-registry-protocol#download-source-code-for-a-specific-module-version
func (gh *GitHubRegistry) Download(namespace, name, provider, version string) string {
	path := gh.Path(provider, name)

	return fmt.Sprintf("git::%s://github.com/%s/%s?ref=v%s", gh.Config.Protocol, namespace, path, version)
}

// Validate GitHub client
func (gh *GitHubRegistry) validate() error {
	if gh.Config.Protocol != "https" && gh.Config.Protocol != "ssh" {
		return errors.New(fmt.Sprintf("Invalid protocol: %s", gh.Config.Protocol))
	}

	return nil
}
