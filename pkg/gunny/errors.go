package gunny

import (
	"fmt"
	"strings"
)

type ErrInvalidDataFormat struct {
	Supplied    string
	ValidValues []string
}

func (e ErrInvalidDataFormat) Error() string {
	return fmt.Sprintf("invalid data format \"%s\"; valid values: %s", e.Supplied, strings.Join(e.ValidValues, ", "))
}

type ErrInvalidDataValueName struct {
	Name string
}

func (e ErrInvalidDataValueName) Error() string {
	return fmt.Sprintf("invalid data value name: %s", e.Name)
}

type ErrValueResolution struct {
	Name  string
	Cause error
}

func (e ErrValueResolution) Error() string {
	return fmt.Sprintf("failed to resolve value for \"%s\": %s", e.Name, e.Cause)
}

func (e ErrValueResolution) Unwrap() error {
	return e.Cause
}

type ErrNamedValuePairParsing struct {
	Cause error
}

func (e ErrNamedValuePairParsing) Error() string {
	return fmt.Sprintf("failed to parse named value pair: %s", e.Cause)
}

func (e ErrNamedValuePairParsing) Unwrap() error {
	return e.Cause
}

type ErrUnrecognizedDataFormat struct {
	Format string
}

func (e ErrUnrecognizedDataFormat) Error() string {
	return fmt.Sprintf("unrecognized data format: %s", e.Format)
}

type ErrTemplateRead struct {
	Cause error
}

func (e ErrTemplateRead) Error() string {
	return fmt.Sprintf("failed to read template; %s", e.Cause)
}

func (e ErrTemplateRead) Unwrap() error {
	return e.Cause
}

type ErrInvalidPipelineConfigFormat struct {
	Supplied    string
	ValidValues []string
}

func (e ErrInvalidPipelineConfigFormat) Error() string {
	return fmt.Sprintf("invalid pipeline configuration format \"%s\"; valid values: %s", e.Supplied, strings.Join(e.ValidValues, ", "))
}

type ErrMissingEnvVar struct {
	Name string
}

func (e ErrMissingEnvVar) Error() string {
	return fmt.Sprintf("missing required environment variable: %s", e.Name)
}

var ErrNoDataSources = fmt.Errorf("pipeline must have at least one data source configured")
