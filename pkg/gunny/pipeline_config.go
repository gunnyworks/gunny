package gunny

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/samber/lo"
)

// PipelineConfigVersion determines which versions of the pipeline
// configuration file format are supported.
type PipelineConfigVersion string

const (
	ConfigVersion1 PipelineConfigVersion = "v1"
)

var (
	supportedPipelineConfigVersions = []PipelineConfigVersion{
		ConfigVersion1,
	}

	validPipelineConfigFormats = []DataFormat{
		DataFormatJSON,
		DataFormatYAML,
	}
)

type Validator interface {
	Validate() error
}

// PipelineConfig encapsulates configuration, from which a Gunny [Pipeline] can
// be constructed.
type PipelineConfig struct {
	// Version of the configuration file. Default: v1.
	Version PipelineConfigVersion `json:"version" yaml:"version"`

	// DataSources lists all the sources of data from which Gunny will attempt
	// to obtain/resolve values associated with named data. The order in which
	// they are specified determines precedence if any duplicate named values
	// are detected.
	DataSources []*DataSourceConfig `json:"data_sources,omitempty" yaml:"data-sources,omitempty"`

	// Renderer configuration, which determines which template engine to use
	// and its associated configuration options.
	Renderer *RendererConfig `json:"renderer,omitempty" yaml:"renderer,omitempty"`

	// OutputBasePath defines the base path for all file-based outputs whose
	// filenames are not specified as being absolute. Defaults to the current
	// working directory.
	OutputBasePath string `json:"output_base_path,omitempty" yaml:"output-base-path,omitempty"`

	// Outputs determines how and where Gunny renders/sends its output. If nil,
	// defaults to stdout output.
	Outputs []*OutputConfig `json:"outputs,omitempty" yaml:"outputs,omitempty"`
}

// ReadPipelineConfig attempts to parse data from the given reader into a
// [PipelineConfig].
func ReadPipelineConfig(r io.Reader, format DataFormat) (*PipelineConfig, error) {
	var config PipelineConfig
	switch format {
	case DataFormatJSON:
		if err := json.NewDecoder(r).Decode(&config); err != nil {
			return nil, err
		}
	case DataFormatYAML:
		if err := yaml.NewDecoder(r).Decode(&config); err != nil {
			return nil, err
		}
	default:
		return nil, InvalidPipelineConfigFormatError{
			Supplied: string(format),
			ValidValues: lo.Map(validPipelineConfigFormats, func(format DataFormat, _ int) string {
				return string(format)
			}),
		}
	}
	// Populate the configuration with default values for anything that wasn't
	// previously set.
	config.PopulateWithDefaults()
	return &config, nil
}

// ReadPipelineConfigFromFile is similar to [ReadPipelineConfig], but it
// handles opening and closing of the file with the given name. Configuration
// file format is inferred from the file's extension.
func ReadPipelineConfigFromFile(filename string, overrideConfigFormat *DataFormat) (*PipelineConfig, error) {
	var configFormat DataFormat
	var err error

	if overrideConfigFormat != nil {
		configFormat = *overrideConfigFormat
	} else {
		configFormat, err = GetFileDataFormatFromFilename(filename)
		if err != nil {
			return nil, err
		}
		switch configFormat {
		case DataFormatJSON:
		case DataFormatYAML:
		default:
			return nil, InvalidPipelineConfigFormatError{
				Supplied: filepath.Ext(filename),
				ValidValues: lo.Map(validPipelineConfigFormats, func(format DataFormat, _ int) string {
					return string(format)
				}),
			}
		}
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	return ReadPipelineConfig(f, configFormat)
}

func NewPipelineConfigWithDefaults() *PipelineConfig {
	config := &PipelineConfig{}
	config.PopulateWithDefaults()
	return config
}

func (c *PipelineConfig) PopulateWithDefaults() {
	if len(c.Version) == 0 {
		c.Version = ConfigVersion1
	}
	if len(c.DataSources) == 0 {
		c.DataSources = []*DataSourceConfig{
			{
				Type: DataSourceStdin,
				Config: &StdinDataSourceConfig{
					Format: DataFormatJSON,
				},
			},
		}
	}
	if c.Renderer == nil {
		c.Renderer = NewRendererConfigWithDefaults()
	}
	if len(c.OutputBasePath) == 0 {
		c.OutputBasePath = "."
	}
	if len(c.Outputs) == 0 {
		c.Outputs = []*OutputConfig{NewOutputConfigWithDefaults()}
	}
}

// Validate that the pipeline configuration does not have any errors.
func (c *PipelineConfig) Validate() error {
	if c.Version != ConfigVersion1 {
		return UnsupportedVersionError{
			Supplied: string(c.Version),
			SupportedVersions: lo.Map(supportedPipelineConfigVersions, func(v PipelineConfigVersion, _ int) string {
				return string(v)
			}),
		}
	}
	if len(c.DataSources) == 0 {
		return ErrMissingDataSource
	}
	if c.Renderer == nil {
		return ErrMissingRendererConfig
	}
	if len(c.Outputs) == 0 {
		return ErrMissingOutputConfig
	}

	for _, dataSource := range c.DataSources {
		if err := dataSource.Validate(); err != nil {
			return err
		}
	}
	if err := c.Renderer.Validate(); err != nil {
		return err
	}
	for _, output := range c.Outputs {
		if err := output.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// SetStdinDataFormat allows one to override the data format associated with
// any stdin data sources.
func (c *PipelineConfig) SetStdinDataFormat(dataFormat DataFormat) {
	for _, dataSource := range c.DataSources {
		if dataSource.Type == DataSourceStdin {
			switch config := dataSource.Config.(type) {
			case *StdinDataSourceConfig:
				config.Format = dataFormat
			}
		}
	}
}

// SetTemplateContent allows one to override the template content associated
// with the renderer.
func (c *PipelineConfig) SetTemplateContent(templateContent string) {
	if c.Renderer == nil {
		c.Renderer = NewRendererConfigWithDefaults()
	}
	c.Renderer.SetTemplateContent(templateContent)
}

// SetTemplateFile allows one to override the template file associated with the
// renderer.
func (c *PipelineConfig) SetTemplateFile(templateFilePath string) {
	if c.Renderer == nil {
		c.Renderer = NewRendererConfigWithDefaults()
	}
	c.Renderer.SetTemplateFile(templateFilePath)
}

// DataSourceType is an enum type referring to the types of data sources from
// which Gunny can obtain data.
type DataSourceType string

const (
	DataSourceUnknown DataSourceType = "unknown"
	DataSourceEnvVars DataSourceType = "env-vars" // Environment variables
	DataSourceCLIArgs DataSourceType = "cli-args" // Command line arguments
	DataSourceStdin   DataSourceType = "stdin"    // Stdin
	DataSourceFile    DataSourceType = "file"     // One or more files (type determined by file extension)
	DataSourceMap     DataSourceType = "map"      // An in-memory map of values
)

var (
	validDataSourceTypeValues = []DataSourceType{
		DataSourceEnvVars,
		DataSourceCLIArgs,
		DataSourceStdin,
		DataSourceFile,
		DataSourceMap,
	}
	validDataSourceTypeStringValues = lo.Map(validDataSourceTypeValues, func(t DataSourceType, _ int) string {
		return string(t)
	})
)

func (t DataSourceType) Validate() error {
	if !lo.Contains(validDataSourceTypeValues, t) {
		return InvalidDataSourceTypeError{
			Supplied:    string(t),
			ValidValues: validDataSourceTypeStringValues,
		}
	}
	return nil
}

type DataSourceConfig struct {
	Type   DataSourceType `json:"type" yaml:"type"`
	Config any            `json:"config,omitempty" yaml:"config, omitempty"`
}

var (
	_ json.Unmarshaler      = (*DataSourceConfig)(nil)
	_ yaml.BytesUnmarshaler = (*DataSourceConfig)(nil)
)

// UnmarshalJSON implements json.Unmarshaler.
func (c *DataSourceConfig) UnmarshalJSON(data []byte) error {
	type proxyType struct {
		Type   DataSourceType  `json:"type"`
		Config json.RawMessage `json:"config,omitempty"`
	}
	var p proxyType
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	c.Type = p.Type
	if len(p.Config) == 0 {
		return nil
	}
	return c.unmarshalConfig(p.Config, json.Unmarshal)
}

// UnmarshalYAML implements yaml.BytesUnmarshaler.
func (c *DataSourceConfig) UnmarshalYAML(data []byte) error {
	type proxyType struct {
		Type   DataSourceType `yaml:"type"`
		Config ast.Node       `yaml:"config,omitempty"`
	}
	var p proxyType
	if err := yaml.Unmarshal(data, &p); err != nil {
		return err
	}
	c.Type = p.Type
	if p.Config == nil {
		return nil
	}
	yamlConfig, err := p.Config.MarshalYAML()
	if err != nil {
		return err
	}
	return c.unmarshalConfig(yamlConfig, yaml.Unmarshal)
}

func (c *DataSourceConfig) unmarshalConfig(data []byte, unmarshalFn func([]byte, any) error) error {
	switch c.Type {
	case DataSourceEnvVars:
		var cfg EnvVarsDataSourceConfig
		if err := unmarshalFn(data, &cfg); err != nil {
			return err
		}
		c.Config = &cfg
	case DataSourceCLIArgs:
		var cfg CLIArgsDataSourceConfig
		if err := unmarshalFn(data, &cfg); err != nil {
			return err
		}
		c.Config = &cfg
	case DataSourceStdin:
		var cfg StdinDataSourceConfig
		if err := unmarshalFn(data, &cfg); err != nil {
			return err
		}
		c.Config = &cfg
	case DataSourceFile:
		var cfg FileDataSourceConfig
		if err := unmarshalFn(data, &cfg); err != nil {
			return err
		}
		c.Config = &cfg
	case DataSourceMap:
		cfg := make(map[string]any)
		if err := unmarshalFn(data, &cfg); err != nil {
			return err
		}
		c.Config = cfg
	default:
		return InvalidDataSourceTypeError{
			Supplied:    string(c.Type),
			ValidValues: validDataSourceTypeStringValues,
		}
	}
	return nil
}

func (c *DataSourceConfig) Validate() error {
	if err := c.Type.Validate(); err != nil {
		return err
	}
	// At present, no data source configurations can be nil.
	if c.Config == nil {
		return InvalidConfigError{
			Type:  c.Type,
			Cause: ErrMissingConfig,
		}
	}
	validator, canValidate := c.Config.(Validator)
	if canValidate {
		return validator.Validate()
	}
	return nil
}

type CLIArgsDataSourceConfig struct {
	// Expected must contain a list of expected CLI argument data value names.
	// Only these named values will be made available during rendering, and an
	// error will be produced if one of them is missing.
	Expected []string `json:"expected,omitempty" yaml:"expected,omitempty"`
	// Optional defines a list of optional CLI argument data values that, if
	// present, will be made available to the rendering pipeline. If any of
	// them are missing, however, no error will be produced.
	Optional []string `json:"optional,omitempty" yaml:"optional,omitempty"`
}

func (c *CLIArgsDataSourceConfig) Validate() error {
	if c == nil {
		return InvalidConfigError{Type: DataSourceCLIArgs, Cause: ErrMissingConfig}
	}
	if len(c.Expected) == 0 && len(c.Optional) == 0 {
		return InvalidConfigError{Type: DataSourceCLIArgs, Cause: ErrRedundantCLIArgsDataSourceConfig}
	}
	return nil
}

type EnvVarsDataSourceConfig struct {
	// Expected must contain a list of expected environment variables. Only
	// these environment variables will be made available during rendering, and
	// an error will be produced if one of them is missing.
	Expected []string `json:"expected,omitempty" yaml:"expected,omitempty"`
	// Optional defines a list of optional environment variables that, if
	// present, will be made available to the rendering pipeline. If any of
	// them are missing, however, no error will be produced.
	Optional []string `json:"optional,omitempty" yaml:"optional,omitempty"`
}

func (c *EnvVarsDataSourceConfig) Validate() error {
	if c == nil {
		return InvalidConfigError{Type: DataSourceEnvVars, Cause: ErrMissingConfig}
	}
	if len(c.Expected) == 0 && len(c.Optional) == 0 {
		return InvalidConfigError{Type: DataSourceEnvVars, Cause: ErrRedundantEnvVarsDataSourceConfig}
	}
	return nil
}

type StdinDataSourceConfig struct {
	// Format defines how Gunny deserializes data from stdin.
	Format DataFormat `json:"format" yaml:"format"`
}

func (c *StdinDataSourceConfig) Validate() error {
	if c == nil {
		return InvalidConfigError{Type: DataSourceStdin, Cause: ErrMissingConfig}
	}
	if err := c.Format.Validate(); err != nil {
		return InvalidConfigError{Type: DataSourceStdin, Cause: err}
	}
	return nil
}

type FileDataSourceConfig struct {
	// Path to the file from which to read data. Can be a glob, in which case
	// multiple files will potentially be matched. When this is the case, data
	// values will be named according to the file name (without extensions,
	// slugified).
	Path string `json:"path" yaml:"path"`
	// Format overrides Gunny's extension-based autodetection for file types.
	// If multiple files are matched and this is specified, all of them will be
	// matched according to this format.
	Format *DataFormat `json:"format,omitempty" yaml:"format,omitempty"`
}

func (c *FileDataSourceConfig) Validate() error {
	if c == nil {
		return InvalidConfigError{Type: DataSourceFile, Cause: ErrMissingConfig}
	}
	if len(c.Path) == 0 {
		return InvalidConfigError{Type: DataSourceFile, Cause: ErrMissingFilePath}
	}
	return nil
}

// RendererType defines an enum type of possible renderers.
type RendererType string

const (
	RendererUnknown  RendererType = "unknown"
	RendererMustache RendererType = "mustache"
)

var (
	validRendererTypes = []RendererType{
		RendererMustache,
	}
	validRendererTypeStrings = lo.Map(validRendererTypes, func(t RendererType, _ int) string {
		return string(t)
	})
)

type RendererConfig struct {
	Type   RendererType `json:"type" yaml:"type"`
	Config any          `json:"config,omitempty" yaml:"config,omitempty"`
}

var (
	_ json.Unmarshaler      = (*RendererConfig)(nil)
	_ yaml.BytesUnmarshaler = (*RendererConfig)(nil)
)

func NewRendererConfigWithDefaults() *RendererConfig {
	return &RendererConfig{
		Type:   RendererMustache,
		Config: &MustacheRendererConfig{},
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *RendererConfig) UnmarshalJSON(data []byte) error {
	type proxyType struct {
		Type   RendererType    `json:"type"`
		Config json.RawMessage `json:"config,omitempty"`
	}
	var p proxyType
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	c.Type = p.Type
	if len(p.Config) == 0 {
		return nil
	}
	return c.unmarshalConfig(p.Config, json.Unmarshal)
}

// UnmarshalYAML implements yaml.BytesUnmarshaler.
func (c *RendererConfig) UnmarshalYAML(data []byte) error {
	type proxyType struct {
		Type   RendererType `yaml:"type"`
		Config ast.Node     `yaml:"config,omitempty"`
	}
	var p proxyType
	if err := yaml.Unmarshal(data, &p); err != nil {
		return err
	}
	c.Type = p.Type
	if p.Config == nil {
		return nil
	}
	yamlConfig, err := p.Config.MarshalYAML()
	if err != nil {
		return err
	}
	return c.unmarshalConfig(yamlConfig, yaml.Unmarshal)
}

func (c *RendererConfig) unmarshalConfig(data []byte, unmarshalFn func([]byte, any) error) error {
	switch c.Type {
	case RendererMustache:
		var cfg MustacheRendererConfig
		if err := unmarshalFn(data, &cfg); err != nil {
			return err
		}
		c.Config = &cfg
		return nil
	}
	return InvalidRendererTypeError{
		Supplied:    string(c.Type),
		ValidValues: validRendererTypeStrings,
	}
}

func (c *RendererConfig) Validate() error {
	switch c.Type {
	case RendererMustache:
		if c.Config == nil {
			return ErrMissingRendererConfig
		}
	}
	validator, canValidate := c.Config.(Validator)
	if canValidate {
		return validator.Validate()
	}
	return nil
}

func (c *RendererConfig) SetTemplateContent(templateContent string) {
	switch rendererConfig := c.Config.(type) {
	case *MustacheRendererConfig:
		rendererConfig.Template = &templateContent
	}
}

func (c *RendererConfig) SetTemplateFile(templateFilePath string) {
	switch rendererConfig := c.Config.(type) {
	case *MustacheRendererConfig:
		rendererConfig.TemplateFile = &templateFilePath
	}
}

// MustacheRendererConfig encapsulates all configuration for rendering of data
// via Mustache templates.
//
// TODO: Add support for Mustache partials.
type MustacheRendererConfig struct {
	// Template gives the option to provide a raw template directly from the
	// configuration.
	Template *string `json:"template,omitempty" yaml:"template,omitempty"`
	// TemplateFile allows one to specify a file from which to load a template.
	TemplateFile *string `json:"template_file,omitempty" yaml:"template-file,omitempty"`
}

func (c *MustacheRendererConfig) Validate() error {
	template := ""
	templateFile := ""
	if c.Template != nil {
		template = *c.Template
	}
	if c.TemplateFile != nil {
		templateFile = *c.TemplateFile
	}
	if len(template) == 0 && len(templateFile) == 0 {
		return ErrMissingTemplate
	}
	if len(template) > 0 && len(templateFile) > 0 {
		return ErrTooManyTemplates
	}
	return nil
}

// OutputType defines an enum type of possible output options.
type OutputType string

const (
	OutputUnknown OutputType = "unknown"
	OutputStdout  OutputType = "stdout" // Output to stdout
	OutputFile    OutputType = "file"   // Output to one file
)

var (
	validOutputTypes = []OutputType{
		OutputStdout,
		OutputFile,
	}

	validOutputTypeStrings = lo.Map(validOutputTypes, func(t OutputType, _ int) string {
		return string(t)
	})
)

type OutputConfig struct {
	Type   OutputType `json:"type" yaml:"type"`
	Config any        `json:"config,omitempty" yaml:"config,omitempty"`
}

var (
	_ json.Unmarshaler      = (*OutputConfig)(nil)
	_ yaml.BytesUnmarshaler = (*OutputConfig)(nil)
)

func NewOutputConfigWithDefaults() *OutputConfig {
	return &OutputConfig{
		Type: OutputStdout,
	}
}

func (c *OutputConfig) Validate() error {
	validator, canValidate := c.Config.(Validator)
	if canValidate {
		return validator.Validate()
	}
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *OutputConfig) UnmarshalJSON(data []byte) error {
	type proxyType struct {
		Type   OutputType      `json:"type"`
		Config json.RawMessage `json:"config,omitempty"`
	}
	var p proxyType
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	c.Type = p.Type
	if len(p.Config) == 0 {
		return nil
	}
	return c.unmarshalConfig(p.Config, json.Unmarshal)
}

// UnmarshalYAML implements yaml.BytesUnmarshaler.
func (c *OutputConfig) UnmarshalYAML(data []byte) error {
	type proxyType struct {
		Type   OutputType `yaml:"type"`
		Config ast.Node   `yaml:"config,omitempty"`
	}
	var p proxyType
	if err := yaml.Unmarshal(data, &p); err != nil {
		return err
	}
	c.Type = p.Type
	if p.Config == nil {
		return nil
	}
	yamlConfig, err := p.Config.MarshalYAML()
	if err != nil {
		return err
	}
	return c.unmarshalConfig(yamlConfig, yaml.Unmarshal)
}

func (c *OutputConfig) unmarshalConfig(data []byte, unmarshalerFn func([]byte, any) error) error {
	switch c.Type {
	case OutputStdout:
		// No configuration needed
		return nil
	case OutputFile:
		var cfg FileOutputConfig
		if err := unmarshalerFn(data, &cfg); err != nil {
			return err
		}
		c.Config = &cfg
		return nil
	}
	return InvalidOutputTypeError{
		Supplied:    string(c.Type),
		ValidValues: validOutputTypeStrings,
	}
}

type FileOutputConfig struct {
	// Path to the file to which to write the output.
	Path string `json:"path" yaml:"path"`
}
