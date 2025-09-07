package gunny

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/samber/lo"
)

type InvalidDataFormatError struct {
	Supplied    string
	ValidValues []string
}

func (e InvalidDataFormatError) Error() string {
	return fmt.Sprintf("invalid data format \"%s\"; valid values: %s", e.Supplied, strings.Join(e.ValidValues, ", "))
}

type InvalidDataValueNameError struct {
	Name string
}

func (e InvalidDataValueNameError) Error() string {
	return fmt.Sprintf("invalid data value name: %s", e.Name)
}

type ValueResolutionError struct {
	Name  string
	Cause error
}

func (e ValueResolutionError) Error() string {
	return fmt.Sprintf("failed to resolve value for \"%s\": %s", e.Name, e.Cause)
}

func (e ValueResolutionError) Unwrap() error {
	return e.Cause
}

type NamedValuePairParsingError struct {
	Cause error
}

func (e NamedValuePairParsingError) Error() string {
	return fmt.Sprintf("failed to parse named value pair: %s", e.Cause)
}

func (e NamedValuePairParsingError) Unwrap() error {
	return e.Cause
}

type UnrecognizedDataFormatError struct {
	Format string
}

func (e UnrecognizedDataFormatError) Error() string {
	return fmt.Sprintf("unrecognized data format: %s", e.Format)
}

type TemplateReadError struct {
	Cause error
}

func (e TemplateReadError) Error() string {
	return fmt.Sprintf("failed to read template; %s", e.Cause)
}

func (e TemplateReadError) Unwrap() error {
	return e.Cause
}

type InvalidPipelineConfigFormatError struct {
	Supplied    string
	ValidValues []string
}

func (e InvalidPipelineConfigFormatError) Error() string {
	return fmt.Sprintf("invalid pipeline configuration format \"%s\"; valid values: %s", e.Supplied, strings.Join(e.ValidValues, ", "))
}

type MissingEnvVarError struct {
	Name string
}

func (e MissingEnvVarError) Error() string {
	return fmt.Sprintf("missing required environment variable: %s", e.Name)
}

type MissingCLIArgError struct {
	Name string
}

func (e MissingCLIArgError) Error() string {
	return fmt.Sprintf("missing required command line argument: %s", e.Name)
}

type InvalidDataSourceTypeError struct {
	Supplied    string
	ValidValues []string
}

func (e InvalidDataSourceTypeError) Error() string {
	return fmt.Sprintf("invalid data source type \"%s\"; valid values: %s", e.Supplied, strings.Join(e.ValidValues, ", "))
}

type InvalidConfigTypeError struct {
	Expected any
	Actual   any
}

func (e InvalidConfigTypeError) Error() string {
	expectedType := "nil"
	actualType := "nil"
	if e.Expected != nil {
		expectedType = reflect.TypeOf(e.Expected).Name()
	}
	if e.Actual != nil {
		actualType = reflect.TypeOf(e.Actual).Name()
	}
	return fmt.Sprintf("expected configuration of type %s, but got %s", expectedType, actualType)
}

type UnrecognizedFileExtError struct {
	Ext           string
	SupportedExts []string
}

func (e UnrecognizedFileExtError) Error() string {
	return fmt.Sprintf("unrecognized file extension \"%s\"; supported file extensions: %s", e.Ext, strings.Join(e.SupportedExts, ", "))
}

type InvalidRendererTypeError struct {
	Supplied    string
	ValidValues []string
}

func (e InvalidRendererTypeError) Error() string {
	return fmt.Sprintf("invalid renderer type \"%s\"; valid values: %s", e.Supplied, strings.Join(e.ValidValues, ", "))
}

type InvalidOutputTypeError struct {
	Supplied    string
	ValidValues []string
}

func (e InvalidOutputTypeError) Error() string {
	return fmt.Sprintf("invalid output type \"%s\"; valid values: %s", e.Supplied, strings.Join(e.ValidValues, ", "))
}

type UnsupportedVersionError struct {
	Supplied          string
	SupportedVersions []string
}

func (e UnsupportedVersionError) Error() string {
	return fmt.Sprintf("unsupported pipeline configuration version \"%s\"; supported versions: %s", e.Supplied, strings.Join(e.SupportedVersions, ", "))
}

type InvalidConfigError struct {
	Type  DataSourceType
	Cause error
}

func (e InvalidConfigError) Error() string {
	return fmt.Sprintf("invalid %s configuration: %s", e.Type, e.Cause)
}

func (e InvalidConfigError) Unwrap() error {
	return e.Cause
}

type MultiWriterCloseFailedError struct {
	Causes []error
}

func (e MultiWriterCloseFailedError) Error() string {
	return fmt.Sprintf(
		"failed to close one or more writers: %s",
		strings.Join(lo.Map(e.Causes, func(c error, _ int) string {
			return c.Error()
		}), ", "),
	)
}

var (
	ErrMissingTemplate                  = errors.New("no template or template file supplied")
	ErrTooManyTemplates                 = errors.New("too many templates - either the template or template file can be specified, but not both")
	ErrMissingDataSource                = errors.New("no data source(s) specified")
	ErrMissingConfig                    = errors.New("missing configuration")
	ErrMissingFilePath                  = errors.New("missing file path")
	ErrMissingRendererConfig            = errors.New("no renderer configuration specified")
	ErrMissingOutputConfig              = errors.New("no output configuration specified")
	ErrRedundantCLIArgsDataSourceConfig = errors.New("redundant (empty) CLI arguments data source configuration")
	ErrRedundantEnvVarsDataSourceConfig = errors.New("redundant (empty) environment variables data source configuration")
)
