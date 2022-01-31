package moduleregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleRegistry(t *testing.T) {
	t.Skip()
	t.Run("GitHub registry", func(t *testing.T) {
		expected := NewGitHubClient("../../test/github-secure.json")
		actual := NewModuleRegistry("github")

		assert.Equal(t, expected, actual, "Registry type should be of type github")
	})
}
