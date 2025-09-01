package gunny

import (
	"encoding/json"
	"io"

	"github.com/goccy/go-yaml"
	"github.com/samber/lo"
)

// DataFormat defines the format in which data is supplied, which gives an
// indication as to which parser Gunny should use.
type DataFormat string

const (
	DataFormatUnrecognized DataFormat = "unrecognized"
	DataFormatJSON         DataFormat = "json"
	DataFormatYAML         DataFormat = "yaml"
)

var validDataFormatValues = []DataFormat{
	DataFormatJSON,
	DataFormatYAML,
}

// DataFormatFromString attempts to convert the given string to a valid
// [DataFormat] value. Returns an error if the supplied value is not
// recognized.
func DataFormatFromString(s string) (DataFormat, error) {
	switch s {
	case string(DataFormatJSON):
		return DataFormatJSON, nil
	case string(DataFormatYAML):
		return DataFormatYAML, nil
	default:
		return DataFormatUnrecognized, ErrInvalidDataFormat{
			Supplied: s,
			ValidValues: lo.Map(validDataFormatValues, func(v DataFormat, _ int) string {
				return string(v)
			}),
		}
	}
}

// NewDataResolverFromReaderattempts to parse data from the given reader into
// an in-memory data map.
func NewDataResolverFromReader(reader io.Reader, format DataFormat) (DataResolverMap, error) {
	switch format {
	case DataFormatJSON:
		return newJSONDataResolverFromReader(reader)
	case DataFormatYAML:
		return newYAMLDataResolverFromReader(reader)
	default:
		return nil, ErrUnrecognizedDataFormat{
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
