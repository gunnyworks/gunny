package gunny_test

import (
	"testing"

	"github.com/gunnyworks/gunny/pkg/gunny"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNameValuePairStringParsing(t *testing.T) {
	testCases := []struct {
		desc         string
		args         []string
		expectErr    bool
		expectedData map[string]any
	}{
		{
			desc: "trivial string name/value pair",
			args: []string{"name=Michael"},
			expectedData: map[string]any{
				"name": "Michael",
			},
		},
		{
			desc: "name/value pair with digits in name",
			args: []string{"somevar02=abcd"},
			expectedData: map[string]any{
				"somevar02": "abcd",
			},
		},
		{
			desc: "name/value pair with numerical value",
			args: []string{"var=12345"},
			expectedData: map[string]any{
				"var": "12345",
			},
		},
		{
			desc: "snake_case name",
			args: []string{"snake_case_var=Some Value"},
			expectedData: map[string]any{
				"snake_case_var": "Some Value",
			},
		},
		{
			desc: "empty value",
			args: []string{"empty_var="},
			expectedData: map[string]any{
				"empty_var": "",
			},
		},
		{
			desc:      "cannot have spaces in name",
			args:      []string{"some var=some value"},
			expectErr: true,
		},
		{
			desc:      "cannot have special characters in name",
			args:      []string{"var%name=some value"},
			expectErr: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			resolver, err := gunny.NewNameValuePairsDataResolver(testCase.args)
			if testCase.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			resolved, err := resolver.Resolve(t.Context())
			require.NoError(t, err)

			data, ok := resolved.(map[string]any)
			require.True(t, ok)

			assert.Equal(t, testCase.expectedData, data)
		})
	}
}
