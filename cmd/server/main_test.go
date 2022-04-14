package main

import (
	"errors"
	"testing"

	"github.com/ironhalo/hellas/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("Github: Valid Config", func(t *testing.T) {
		expected := &models.Config{
			ModuleBackend: "github",
			ModuleRegistry: &models.ModuleRegistry{
				InsecureSkipVerify: true,
				Protocol:           "https",
				Prefix:             "terraform",
			},
		}

		data := `
		{
			"moduleBackend": "github",
			"moduleRegistry": {
				"insecureSkipVerify": true,
				"protocol": "https",
				"prefix": "terraform"
			}
		}`

		actual, err := newConfig([]byte(data))

		assert.Equal(t, expected, actual)
		assert.Equal(t, nil, err)
	})

	t.Run("Github: Invalid Protocol", func(t *testing.T) {
		data := `
		{
			"moduleBackend": "github",
			"moduleRegistry": {
				"insecureSkipVerify": true,
				"protocol": "foo",
				"prefix": "terraform"
			}
		}`

		_, err := newConfig([]byte(data))

		assert.Equal(t, errors.New("Invalid protocol: foo"), err)
	})
}
