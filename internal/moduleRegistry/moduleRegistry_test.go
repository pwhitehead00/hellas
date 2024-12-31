package moduleregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModuleRegistry(t *testing.T) {
	testCases := []struct {
		desc  string
		input Config
		err   error
	}{
		{
			desc: "Valid GitHub Config",
			input: Config{
				Registries: Registries{
					Github: GithubConfig{
						Protocol:           "https",
						InsecureSkipVerify: true,
						Enabled:            true,
					},
				},
			},
			err: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := NewModuleRegistry(tC.input)
			assert.Equal(t, tC.err, err)
		})
	}
}
