package moduleregistry

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/v40/github"
	"github.com/ironhalo/hellas/models"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	Client *github.Client
	Config *GitHubConfig
}

type GitHubConfig struct {
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
	Protocol           string `json:"protocol"`
	Prefix             string `json:"prefix"`
}

func initConfig() (*GitHubConfig, error) {
	var config GitHubConfig
	f, ok := os.LookupEnv("CONFIG")
	if !ok {
		log.Fatal("Failed to load config")
	}

	file, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	if config.Protocol != "https" && config.Protocol != "ssh" {
		log.Fatalf("Invalid protocol: %s", config.Protocol)
	}

	return &config, nil
}

func NewGitHubClient() Registry {
	config, err := initConfig()
	if err != nil {
		log.Fatalf("Failed to initialize github config: %s", err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify},
	}

	sslcli := &http.Client{Transport: tr}
	ctx := context.TODO()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		log.Println("No token provided, using unauthenticated GitHub client")
		tc := oauth2.NewClient(ctx, nil)
		client := github.NewClient(tc)

		return &GitHubClient{
			Client: client,
			Config: config,
		}
	}

	log.Println("Token found, using authenticated GitHub client")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(token)},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubClient{
		Client: client,
		Config: config,
	}
}

func repo(prefix, provider, name string) string {
	if prefix == "" {
		return fmt.Sprintf("%s-%s", provider, name)
	}

	return fmt.Sprintf("%s-%s-%s", prefix, provider, name)
}

func (gh *GitHubClient) GetVersions(namespace, name, provider string) ([]string, error) {
	var allTags []*github.RepositoryTag
	var versions []string

	opt := &github.ListOptions{
		PerPage: 100,
	}

	repo := repo(gh.Config.Prefix, provider, name)

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

func (gh *GitHubClient) Versions(namespace, name, provider string, version []string) models.ModuleVersions {
	var m models.ModuleVersions
	var mv []*models.ModuleVersion

	repo := repo(gh.Config.Prefix, provider, name)

	for _, t := range version {
		o := models.ModuleVersion{
			Version: t,
		}
		mv = append(mv, &o)
	}

	mpv := models.ModuleProviderVersions{
		Source:   fmt.Sprintf("%s/%s", namespace, repo),
		Versions: mv,
	}

	m.Modules = append(m.Modules, &mpv)

	return m
}

func (gh *GitHubClient) Download(namespace, name, provider, version string) string {
	if gh.Config.Prefix == "" {
		return fmt.Sprintf("git::%s://github.com/%s/%s-%s?ref=v%s", gh.Config.Protocol, namespace, provider, name, version)
	}
	return fmt.Sprintf("git::%s://github.com/%s/%s-%s-%s?ref=v%s", gh.Config.Protocol, namespace, gh.Config.Prefix, provider, name, version)

}
