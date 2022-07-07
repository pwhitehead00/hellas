package moduleregistry

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/xanzy/go-gitlab"
)

type gitLabRegistry struct {
	Client *gitlab.Client
	Config *gitLabConfig
}

type gitLabConfig struct {
	// Accept self signed certs
	InsecureSkipVerify bool `json:"insecureSkipVerify"`
	// The protocol used when checking out terraform modules.  Can be either
	// https or ssh
	Protocol string `json:"protocol"`
	// Set a custom URL for the gitlab instance
	BaseURL string `json:"baseURL"`
	// Set the parent groups for the Terraform module
	Groups string `json:"groups"`
	// Derive hostname from BaseURL
	HostName string
}

func newGitLabConfig(file []byte) (*gitLabConfig, error) {
	var config gitLabConfig

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func newGitLabRegistry(config *gitLabConfig) Registry {
	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		log.Println("No token provided, using unauthenticated GitLab client")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify},
	}
	c := &http.Client{Transport: tr}

	client, err := gitlab.NewClient(token, gitlab.WithHTTPClient(c), gitlab.WithBaseURL(config.BaseURL))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return &gitLabRegistry{
		Client: client,
		Config: config,
	}
}

// Helper function to build the GitLab repo path
func (gl *gitLabRegistry) Path(namespace, provider, name string) string {
	if gl.Config.Groups == "" {
		return fmt.Sprintf("%s/%s/%s", namespace, provider, name)
	}

	return fmt.Sprintf("%s/%s/%s/%s", gl.Config.Groups, namespace, provider, name)
}

// List all tags for a GitLab project
func (gl *gitLabRegistry) ListVersions(namespace, name, provider string) ([]string, error) {
	var allTags []*gitlab.Tag
	var versions []string

	opt := &gitlab.ListTagsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
		},
	}

	repo := gl.Path(namespace, provider, name)

	for {
		tags, resp, err := gl.Client.Tags.ListTags(repo, opt)
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
		versions = append(versions, v.Name)
	}
	return versions, nil
}

// Download source code for a specific module version
// See https://www.terraform.io/internals/module-registry-protocol#download-source-code-for-a-specific-module-version
func (gl *gitLabRegistry) Download(namespace string, name string, provider string, version string) string {
	if gl.Config.Groups == "" {
		return fmt.Sprintf("git::%s://gitlab.com/%s/%s/%s?ref=v%s", gl.Config.Protocol, namespace, provider, name, version)
	}
	return fmt.Sprintf("git::%s://gitlab.com/%s/%s/%s/%s?ref=v%s", gl.Config.Protocol, gl.Config.Groups, namespace, provider, name, version)
}

// Validate GitLab config
func (gl *gitLabRegistry) validate() error {
	if gl.Config.Protocol != "https" && gl.Config.Protocol != "ssh" {
		return fmt.Errorf("Invalid protocol: %s", gl.Config.Protocol)
	}

	if gl.Config.BaseURL != "" {
		u, err := url.Parse(gl.Config.BaseURL)
		if err != nil {
			return err
		}

		if u.Scheme != "https" && u.Scheme != "http" {
			return fmt.Errorf("Invalid scheme, only http(s) is supported, got %s", u.Scheme)
		}

		gl.Config.HostName = u.Host
	}

	if gl.Config.BaseURL == "" {
		gl.Config.BaseURL = "https://gitlab.com"
		gl.Config.HostName = "gitlab.com"
	}

	return nil
}
