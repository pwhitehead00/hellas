package moduleregistry

import (
	"testing"

	"github.com/ironhalo/hellas/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestVersions(t *testing.T) {
	t.Run("GitHub Versions", func(t *testing.T) {
		expected := models.ModuleVersions{
			Modules: []*models.ModuleProviderVersions{
				{
					Source: "my-namespace/prefix-happycloud-module",
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
		repo := mr.Repo("happycloud", "module")
		actual := Versions("my-namespace", "module", "happycloud", repo, []string{"1.0.0", "1.0.1"})
		assert.Equal(t, expected, actual, "Versions should be the same")
	})
}
