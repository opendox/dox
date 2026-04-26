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

  @File    : server/README.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-24
  @Modified: 2026-04-26
-->

# Dox Server

`server` is the Web backend runtime for Dox.

The current module contains the CLI entrypoint, shared version command, bootstrap configuration snapshot loading, and the initial server-owned identity setting aggregate. HTTP server startup, database access, logging, security, and EDA integration are intentionally out of scope for this milestone.

## Configuration Bootstrap

`server/internal/bootstrap` can load a startup configuration snapshot through `packages/shared/config`.

The current bootstrap convention is:

- `configs/base.<format>` as the required baseline source;
- `configs/<env>.<format>` as an optional environment override;
- `configs/local.<format>` as an optional local override;
- `DOX_SERVER_` environment variables as optional final overrides.

The bootstrap snapshot currently uses `map[string]any`. Concrete server setting groups belong under `server/internal/setting`, where each configuration group owns its own file. The first group is identity; concrete HTTP, database, cache, logger, security, and IAM setting structs remain out of scope until those runtime resources are introduced.

## Usage

From the repository root:

```bash
go run ./server version
```

From the `server` module:

```bash
go run . version
```

Build the server CLI binary from the repository root:

```bash
go build -o dox-server ./server
```
