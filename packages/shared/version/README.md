<!--
  dox
  Copyright (C) 2026  OpenDox

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program. If not, see <http://www.gnu.org/licenses/>.

  @File    : packages/shared/version/README.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-24
  @Modified: 2026-04-24
-->

# Shared Version Package

`packages/shared/version` provides build and source metadata for Dox backend binaries.

The package is intentionally dependency-free. Build systems should inject metadata through Go `-ldflags -X`, while local development builds keep stable fallback values.

## Build Metadata

Supported build variables:

```text
github.com/opendox/dox/packages/shared/version.buildName
github.com/opendox/dox/packages/shared/version.buildVersion
github.com/opendox/dox/packages/shared/version.buildTime
github.com/opendox/dox/packages/shared/version.buildUser
github.com/opendox/dox/packages/shared/version.buildCGOEnabled
github.com/opendox/dox/packages/shared/version.buildGitCommit
github.com/opendox/dox/packages/shared/version.buildGitBranch
github.com/opendox/dox/packages/shared/version.buildGitTag
github.com/opendox/dox/packages/shared/version.buildGitDirty
```

Example:

```bash
go build -ldflags "\
  -X github.com/opendox/dox/packages/shared/version.buildName=dox-server \
  -X github.com/opendox/dox/packages/shared/version.buildVersion=0.1.0 \
  -X github.com/opendox/dox/packages/shared/version.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -X github.com/opendox/dox/packages/shared/version.buildGitCommit=$(git rev-parse --short HEAD) \
  -X github.com/opendox/dox/packages/shared/version.buildGitBranch=$(git rev-parse --abbrev-ref HEAD) \
  -X github.com/opendox/dox/packages/shared/version.buildGitDirty=$(test -z "$(git status --porcelain)" && echo false || echo true)" \
  ./server/cmd/dox-server
```
