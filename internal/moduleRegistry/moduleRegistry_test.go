package moduleregistry

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleRegistry(t *testing.T) {
	t.Run("GitHub registry", func(t *testing.T) {
		rt := "github"

		config := json.RawMessage(`{"insecureSkipVerify":true,"prefix":"prefix","protocol":"https"}`)
		actual, _ := NewModuleRegistry(&rt, config)

		c := &gitHubConfig{
			InsecureSkipVerify: true,
			Prefix:             "prefix",
			Protocol:           "https",
		}
		expected, _ := NewGitHubRegistry(c)
		assert.Equal(t, expected, actual)
	})

	t.Run("Bad registry", func(t *testing.T) {
		rt := "foobar"

		config := json.RawMessage(`{""}`)
		actual, err := NewModuleRegistry(&rt, config)

		assert.Equal(t, nil, actual)
		assert.Equal(t, errors.New("Unsupported registy type: foobar"), err)
	})
}
