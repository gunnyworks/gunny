package gunny

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMustacheTemplateRendering(t *testing.T) {
	testCases := []struct {
		desc           string
		rawTemplate    string
		data           string
		dataFormat     DataFormat
		expectedOutput string
	}{
		{
			desc:           "trivial example from README",
			rawTemplate:    `Hello {{name}}!`,
			data:           `{"name": "Michael"}`,
			dataFormat:     DataFormatJSON,
			expectedOutput: `Hello Michael!`,
		},
		{
			desc:           "simple example with multiple data types",
			rawTemplate:    `Hello {{name}}! You have {{count}} new message(s)`,
			data:           `{"name": "Michael", "count": 4}`,
			dataFormat:     DataFormatJSON,
			expectedOutput: `Hello Michael! You have 4 new message(s)`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			var actualOutput strings.Builder
			pipeline, err := NewPipeline(
				WithMustacheTemplateFromReader(strings.NewReader(testCase.rawTemplate)),
				WithDataFromReader(strings.NewReader(testCase.data), testCase.dataFormat),
				WithOutputWriter(&actualOutput),
			)
			require.NoError(t, err)

			err = pipeline.Render(t.Context())
			require.NoError(t, err)

			require.Equal(t, testCase.expectedOutput, actualOutput.String())
		})
	}
}

func TestDataResolverMergingAndPrecedence(t *testing.T) {
	testCases := []struct {
		desc           string
		rawTemplate    string
		cliArgs        []string
		data           string
		dataFormat     DataFormat
		expectedOutput string
	}{
		{
			desc:           "trivial example from README",
			rawTemplate:    `Hello {{name}}!`,
			cliArgs:        []string{"name=Gary"},
			data:           `{"name": "Michael"}`,
			dataFormat:     DataFormatJSON,
			expectedOutput: `Hello Gary!`,
		},
		{
			desc:           "simple example with multiple data types",
			rawTemplate:    `Hello {{name}}! You have {{count}} new message(s)`,
			cliArgs:        []string{"count=5"},
			data:           `{"name": "Michael", "count": 4}`,
			dataFormat:     DataFormatJSON,
			expectedOutput: `Hello Michael! You have 5 new message(s)`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			var actualOutput strings.Builder
			pipeline, err := NewPipeline(
				WithMustacheTemplateFromReader(strings.NewReader(testCase.rawTemplate)),
				WithDataFromReader(strings.NewReader(testCase.data), testCase.dataFormat),
				WithDataFromNameValuePairs(testCase.cliArgs),
				WithOutputWriter(&actualOutput),
			)
			require.NoError(t, err)

			err = pipeline.Render(t.Context())
			require.NoError(t, err)

			require.Equal(t, testCase.expectedOutput, actualOutput.String())
		})
	}
}
