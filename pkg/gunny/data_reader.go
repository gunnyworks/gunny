package gunny

import (
	"encoding/json"
	"io"

	"github.com/goccy/go-yaml"
)

type DataFormat string

const (
	DataFormatJSON DataFormat = "json"
	DataFormatYAML DataFormat = "yaml"
)

// NewDataResolverFromReader attempts to parse data from the given reader into
// an in-memory data map.
func NewDataResolverFromReader(reader io.Reader, format DataFormat) (DataResolverMap, error) {
	switch format {
	case DataFormatJSON:
		return newJSONDataResolverFromReader(reader)
	case DataFormatYAML:
		return newYAMLDataResolverFromReader(reader)
	default:
		return nil, UnrecognizedDataFormatError{
			Format: string(format),
		}
	}
}

func newJSONDataResolverFromReader(reader io.Reader) (DataResolverMap, error) {
	var jsonValue map[string]any
	if err := json.NewDecoder(reader).Decode(&jsonValue); err != nil {
		return nil, err
	}
	return NewInMemoryDataResolverMap(jsonValue)
}

func newYAMLDataResolverFromReader(reader io.Reader) (DataResolverMap, error) {
	var yamlValue map[string]any
	if err := yaml.NewDecoder(reader).Decode(&yamlValue); err != nil {
		return nil, err
	}
	return NewInMemoryDataResolverMap(yamlValue)
}
