package gunny_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gunnyworks/gunny/pkg/gunny"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONDataResolver(t *testing.T) {
	testCases := []struct {
		desc         string
		rawJSON      string
		expectedData map[string]any
	}{
		{
			desc:    "trivial string value",
			rawJSON: `{"name": "value"}`,
			expectedData: map[string]any{
				"name": "value",
			},
		},
		{
			desc:    "variables with multiple types",
			rawJSON: `{"string_value": "value", "bool_value": true, "number": 123}`,
			expectedData: map[string]any{
				"string_value": "value",
				"bool_value":   true,
				"number":       123,
			},
		},
		{
			desc: "sub-objects",
			rawJSON: `{
				"some_value": "string value",
				"object": {
					"bool": false,
					"number": 987
				}
			}`,
			expectedData: map[string]any{
				"some_value": "string value",
				"object": map[string]any{
					"bool":   false,
					"number": 987,
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			resolver, err := gunny.NewDataResolverFromReader(strings.NewReader(testCase.rawJSON), gunny.DataFormatJSON)
			require.NoError(t, err)

			resolvedData, err := resolver.Resolve(t.Context())
			require.NoError(t, err)

			resolvedDataMap, ok := resolvedData.(map[string]any)
			require.True(t, ok)

			assertMapDeepEqual(t, testCase.expectedData, resolvedDataMap)
		})
	}
}

func TestYAMLDataResolver(t *testing.T) {
	testCases := []struct {
		desc         string
		rawYAML      string
		expectedData map[string]any
	}{
		{
			desc:    "trivial string value",
			rawYAML: `name: value`,
			expectedData: map[string]any{
				"name": "value",
			},
		},
		{
			desc: "variables with multiple types",
			rawYAML: `string_value: value
bool_value: true
number: 123`,
			expectedData: map[string]any{
				"string_value": "value",
				"bool_value":   true,
				"number":       123,
			},
		},
		{
			desc: "sub-objects",
			rawYAML: `
some_value: string value
object:
  bool: false
  number: 987`,
			expectedData: map[string]any{
				"some_value": "string value",
				"object": map[string]any{
					"bool":   false,
					"number": 987,
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			resolver, err := gunny.NewDataResolverFromReader(strings.NewReader(testCase.rawYAML), gunny.DataFormatYAML)
			require.NoError(t, err)

			resolvedData, err := resolver.Resolve(t.Context())
			require.NoError(t, err)

			resolvedDataMap, ok := resolvedData.(map[string]any)
			require.True(t, ok)

			assertMapDeepEqual(t, testCase.expectedData, resolvedDataMap)
		})
	}
}

func assertMapDeepEqual(t *testing.T, expected map[string]any, actual map[string]any) {
	t.Helper()

	for name, expectedValue := range expected {
		if reflect.ValueOf(expectedValue).Kind() == reflect.Map {
			expectedValueMap, ok := expectedValue.(map[string]any)
			require.True(t, ok)

			actualMap, ok := actual[name].(map[string]any)
			require.True(t, ok)
			assertMapDeepEqual(t, expectedValueMap, actualMap)
			continue
		}
		assert.EqualValues(t, expectedValue, actual[name])
	}
}
