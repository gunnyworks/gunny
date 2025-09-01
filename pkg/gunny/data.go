package gunny

import (
	"context"
	"maps"
	"regexp"
)

var validDataValueNameRegexp = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")

// DataResolver resolves the value of specific data.
type DataResolver interface {
	Resolve(ctx context.Context) (any, error)
}

// DataResolverMap implements a collection of data resolvers where each
// resolver can be referenced by name.
type DataResolverMap map[string]DataResolver

var _ DataResolver = make(DataResolverMap)

// NewInMemoryDataResolverMap creates a [DataResolverMap] in memory from the
// given named map of generic values.
func NewInMemoryDataResolverMap(src map[string]any) (DataResolverMap, error) {
	result := make(DataResolverMap, len(src))
	for name, value := range src {
		if !validDataValueNameRegexp.MatchString(name) {
			return nil, ErrInvalidDataValueName{Name: name}
		}
		result[name] = NewInMemoryDataValue(value)
	}
	return result, nil
}

// Resolve implements [DataResolver].
func (v DataResolverMap) Resolve(ctx context.Context) (any, error) {
	resolvedValues := make(map[string]any, len(v))
	for name, value := range v {
		resolvedValue, err := value.Resolve(ctx)
		if err != nil {
			return nil, ErrValueResolution{
				Name:  name,
				Cause: err,
			}
		}
		resolvedValues[name] = resolvedValue
	}
	return resolvedValues, nil
}

// Merge takes the given other [DataResolverMap] map and merges it into this
// one, returning this map. Values from the supplied (other) map take
// precedence over values in this one.
func (v DataResolverMap) Merge(o DataResolverMap) DataResolverMap {
	maps.Copy(v, o)
	return v
}

// InMemoryDataValue represents a simple in-memory value, where resolution of
// the value produces the value itself.
type InMemoryDataValue[T any] struct {
	value T
}

// NewInMemoryDataValue wraps the given value in an [InMemoryDataValue].
func NewInMemoryDataValue[T any](value T) *InMemoryDataValue[T] {
	return &InMemoryDataValue[T]{
		value: value,
	}
}

// Resolve implements [DataResolver].
func (v *InMemoryDataValue[T]) Resolve(ctx context.Context) (any, error) {
	return v.value, nil
}
