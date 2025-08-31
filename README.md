# Gunny

Fast, simple static text generation from templates and data.

## Usage

### Inline templating and data substitution

```bash
gunny -t 'Hello {{name}}!' -v name=Michael
# Hello Michael!
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
