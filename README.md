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
gunny -t 'Hello {{name}}!' -v name=Michael
# Hello Michael!

echo '{"name": "Gary"}' | gunny -t 'Hello {{name}}!'
```

## Versioning

Gunny follows [semantic versioning](https://semver.org). Prior to v1.0, Gunny's
interfaces (code and CLI) are subject to breaking changes in minor version
releases (e.g. v0.1 to v0.2).

## License

Copyright 2025 Thane Thomson and contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
