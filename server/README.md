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
  @Modified: 2026-04-24
-->

# Dox Server

`server` is the Web backend runtime for Dox.

The current module only contains the CLI entrypoint and shared version command. HTTP server startup, configuration, logging, database access, and EDA integration are intentionally out of scope for this milestone.

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
