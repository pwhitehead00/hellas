package moduleregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateGitHub(t *testing.T) {
	testCases := []struct {
		desc  string
		input githubConfig
		err   error
	}{
		{
			desc: "Successful Github validation",
			input: githubConfig{
				Protocol: "https",
			},
			err: nil,
		},
		{
			desc: "failed Github validation",
			input: githubConfig{
				Protocol: "bad protocol",
			},
			err: invalidProtocol,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.input.validate()
			assert.ErrorIs(t, actual, tC.err)
		})
	}
}
