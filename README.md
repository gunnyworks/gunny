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

## Versioning

Gunny follows [semantic versioning](https://semver.org). Prior to v1.0, Gunny's
interfaces (code and CLI) are subject to breaking changes in minor version
releases (e.g. v0.1 to v0.2).

## License

Copyright 2025 Thane Thomson and contributors

Licensed under a modified version of the Elastic License 2.0. See
[LICENSE](./LICENSE) for details. If you wish to make use of the software for
purposes not covered by the license, please contact the author(s).
