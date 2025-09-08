# Gunny

[![Main Branch Build Status](https://img.shields.io/github/actions/workflow/status/gunnyworks/gunny/go.yaml)](https://github.com/gunnyworks/gunny/actions/workflows/go.yaml?query=branch%3Amain)
[![GitHub Release](https://img.shields.io/github/v/release/gunnyworks/gunny)](https://github.com/gunnyworks/gunny/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/gunnyworks/gunny.svg)](https://pkg.go.dev/github.com/gunnyworks/gunny)

Fast, simple static text generation from templates and data.

## Installation

### Prerequisites

Supported platforms: 

- macOS 
- Linux

### Installing from GitHub release

Visit [Releases](https://github.com/gunnyworks/gunny/releases) and download the
latest version for your platform, installing it somewhere on your system path.

### Installing from source

In order to build from source, you will need [Go v1.25.0](https://go.dev) or
higher installed.

```bash
go install github.com/gunnyworks/gunny/cmd/gunny@latest
```

## Usage

### Inline templating and data substitution

```bash
# Supply template and data via CLI arguments
gunny -t 'Hello {{name}}!' -d name=Michael
# Hello Michael!

# Supply template via CLI argument, but pipe data in in JSON format
echo '{"name": "Gary"}' | gunny -t 'Hello {{name}}!'
# Hello Gary!

# Pipe data in in YAML format
echo 'name: Sarah' | gunny -t 'Hello {{name}}!' --stdin-format yaml
# Hello Sarah!
```

### File-based templates

```bash
echo 'Hello {{name}}!' > ./hello.mustache
gunny --template-file ./hello.mustache -d name=Michael
# Hello Michael!

echo '{"name": "Gary"}' | gunny --template-file ./hello.mustache
# Hello Gary!
```

### Pipeline configuration files

Gunny supports pipeline construction from configuration files. At present, only
JSON- and YAML-based configuration file formats are supported. Following is an
example of a simple configuration file, defined in YAML:

```yaml
# pipeline.yaml

# Configuration file format version
version: v1

# Data sources are processed in sequence, so in cases of data values with
# conflicting names, the later data sources' values take precedence.
data-sources:
  # First, we specify a map of named data values.
  - type: map
    config:
      name: Michael
      count: 4
  # Then, we allow any data from stdin (specified in YAML format) to override
  # our map.
  - type: stdin
    config:
      format: yaml

# Renderer configuration, which controls which template engine we use as well
# as the settings to apply to that template engine.
renderer:
  type: mustache
  config:
    template: "Hello {{name}}! You have {{count}} new messages."

# Gunny supports multiple outputs. By default, if no output is specified,
# stdout will be used.
outputs:
  - type: stdout
```

Then, running Gunny should give:

```bash
gunny -c pipeline.yaml
# Hello Michael! You have 4 new messages.

echo 'name: Sarah' | gunny -c pipeline.yaml
# Hello Sarah! You have 4 new messages.
```

## Pipeline Configuration

Gunny pipeline configuration files can, at present, be specified in either YAML
or JSON format. Following is the high-level format in YAML:

```yaml
# Should always be v1 for now.
version: v1

# Data sources are specified in order of precedence, from lowest to highest.
data-sources:
  # An array of data sources, each of the form:
  - type: <data source type>
    config:
      # Configuration for this particular data source type.

# A single Gunny pipeline only uses a single renderer.
renderer:
  # Which type of renderer should be used?
  type: <renderer type>
  # Renderer configuration
  config:
    # Configuration for this particular type of renderer.

# A list of outputs to which the output data should be written.
outputs:
  - type: <output type>
    config:
      # Configuration for this particular type of output.
```

### Data Sources

Gunny obtains its data from one or more data sources in order to make it
available to the renderer. The following data source types are available:

```yaml
data-sources:
  # Environment variables (an allow-list)
  - type: env-vars
    config:
      # A list of environment variables to make available to the renderer that
      # are EXPECTED to be present at the time of loading the pipeline
      # configuration. If any of these variables are not set, Gunny will throw
      # an error.
      expected:
        - ENVVAR1
        - ENVVAR2
      # A list of OPTIONAL environment variables to make available to the
      # renderer. These work the same as the expected ones, except Gunny does
      # not throw an error if any of them are not specified.
      optional:
        - ENVVAR3
        - ...

  # Command line arguments (specified by way of "-d name=value" flags)
  - type: cli-args
    config:
      # A list of command line argument-specified data values to make available
      # to the renderer that are EXPECTED to be present. If any of these values
      # are not specified via command line arguments, Gunny will throw an
      # error.
      expected:
        - arg1
        - arg2
      optional:
        - arg3
        - ...

  # Read data from stdin, parse and interpret it as named data values
  - type: stdin
    config:
      # Either "json" or "yaml"
      format: json|yaml

  # Read data from a file, parse and interpret it as named data values
  - type: file
    config:
      # Required path to the file to read. If the format field is not specified
      # (see below), the file format will attempt to be determined from this
      # file's extension.
      path: /path/to/file.json
      # OPTIONAL format override. Useful when the file does not have an
      # explicit, or has an unrelated, extension.
      format: json|yaml

  # A simple named data value map, whose names and values are explicitly
  # specified in the pipeline configuration file itself.
  - type: map
    config:
      # Arbitrary key/value pairs (values themselves can be nested objects)
      hello: "world!"
      a_number: 42
```

### Renderers

At present, Gunny only supports [Mustache](https://mustache.github.io/)
templates. More template engines are planned to be supported in future.

```yaml
renderer:
  type: mustache
  # NOTE: Only ONE of either the template or template file can be specified. If
  # both (or neither) are specified, Gunny will fail with an error.
  config:
    # If you want to specify an inline template in the pipeline configuration
    # itself.
    template: "Mustache {{template}}"
    # If you want Gunny to load the Mustache template from a file
    template-file: /path/to/template.mustache
```

### Outputs

Gunny can write the renderer output data to one or more output destinations.

```bash
# A list of outputs to which to send rendered data. Order does not matter, as
# they all receive the same output.
outputs:
  # The default output if no outputs are explicitly specified.
  # Write to stdout. This option has no configuration associated with it.
  - type: stdout

  # Write to a file.
  - type: file
    config:
      # Where to write the output data
      path: /path/to/output/file
```

## Versioning

Gunny follows [semantic versioning](https://semver.org). Prior to v1.0, Gunny's
interfaces (code and CLI) are subject to breaking changes in minor version
releases (e.g. v0.1 to v0.2).

## License

Copyright 2025 Thane Thomson and contributors

Licensed under a modified version of the Elastic License 2.0. See
[LICENSE](./LICENSE) for details. If you wish to make use of the software for
purposes not covered by the license, please contact the author(s).
