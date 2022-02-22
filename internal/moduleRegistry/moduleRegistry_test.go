package moduleregistry

import (
	"testing"

	"github.com/ironhalo/hellas/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestModuleRegistry(t *testing.T) {
	t.Skip()
	mr := models.ModuleRegistry{
		InsecureSkipVerify: true,
		Protocol:           "https",
		Prefix:             "prefix",
	}
	t.Run("GitHub registry", func(t *testing.T) {
		expected := NewGitHubClient(mr)
		actual := NewModuleRegistry("github", mr)

		assert.Equal(t, expected, actual, "Registry type should be of type github")
	})
}
