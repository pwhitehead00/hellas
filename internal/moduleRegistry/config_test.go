package moduleregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateGitHub(t *testing.T) {
	testCases := []struct {
		desc  string
		input GithubConfig
		err   error
	}{
		{
			desc: "Successful Github validation",
			input: GithubConfig{
				Protocol: "https",
			},
			err: nil,
		},
		{
			desc: "failed Github validation",
			input: GithubConfig{
				Protocol: "bad protocol",
			},
			err: invalidProtocol,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.input.Validate()
			assert.ErrorIs(t, actual, tC.err)
		})
	}
}
