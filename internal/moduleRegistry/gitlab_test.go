package moduleregistry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitLabValidation(t *testing.T) {
	t.Run("Gitlab: Invalid Protocol", func(t *testing.T) {
		c := &gitLabConfig{
			Protocol: "foo",
		}

		mr := newGitLabRegistry(c)
		err := mr.validate()

		assert.Equal(t, errors.New("Invalid protocol: foo"), err)
	})

	t.Run("Gitlab: Invalid Scheme", func(t *testing.T) {
		c := &gitLabConfig{
			Protocol: "https",
			BaseURL:  "foo://gitlab.com",
		}

		mr := newGitLabRegistry(c)
		actual := mr.validate()

		assert.Equal(t, errors.New("Invalid scheme, only http(s) is supported, got foo"), actual)
	})
}
