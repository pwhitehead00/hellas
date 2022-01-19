package moduleregistry

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleRegistry(t *testing.T) {
	t.Run("GitHub registry", func(t *testing.T) {
		os.Setenv("CONFIG", "../../test/github-default.json")
		expected := NewGitHubClient()
		actual := NewModuleRegistry("github")

		assert.Equal(t, expected, actual, "Registry type should be of type github")
	})
}
