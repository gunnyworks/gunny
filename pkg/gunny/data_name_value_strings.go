package gunny

import (
	"fmt"
	"regexp"
)

var (
	nameValuePairRegexp     = regexp.MustCompile(`^(?P<name>[a-zA-Z][a-zA-Z0-9_]*)=(?P<value>.*)$`)
	nameValuePairNameIndex  = nameValuePairRegexp.SubexpIndex("name")
	nameValuePairValueIndex = nameValuePairRegexp.SubexpIndex("value")
)

// NewNameValuePairsDataResolver creates a new [DataResolverMap] from the
// given "name=value" pairs (usually supplied via the command line).
func NewNameValuePairsDataResolver(args []string) (DataResolverMap, error) {
	resolvers := make(DataResolverMap, len(args))
	for _, arg := range args {
		name, value, err := parseNameValuePairString(arg)
		if err != nil {
			return nil, err
		}
		resolvers[name] = NewInMemoryDataValue(value)
	}
	return resolvers, nil
}

func parseNameValuePairString(arg string) (string, string, error) {
	matches := nameValuePairRegexp.FindStringSubmatch(arg)
	if len(matches) == 0 {
		return "", "", NamedValuePairParsingError{
			Cause: fmt.Errorf("invalid format for argument \"%s\"", arg),
		}
	}
	name := matches[nameValuePairNameIndex]
	if len(name) == 0 {
		return "", "", NamedValuePairParsingError{
			Cause: fmt.Errorf("failed to parse data value name from argument \"%s\"", arg),
		}
	}
	value := matches[nameValuePairValueIndex]
	return name, value, nil
}
