package gunny

import "fmt"

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
